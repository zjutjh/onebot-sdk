// Package gen 实现 NapCat OpenAPI 到 Go API 绑定的生成器。
package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/pb33f/libopenapi"
	highbase "github.com/pb33f/libopenapi/datamodel/high/base"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	gogen "github.com/pb33f/libopenapi/generator/golang"
	"github.com/pb33f/libopenapi/orderedmap"
)

// GenerateFromFile 从 OpenAPI spec 生成 API 代码。
func GenerateFromFile(specPath string, outDir string) error {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return fmt.Errorf("读取 OpenAPI spec 失败: %w", err)
	}
	doc, err := libopenapi.NewDocument(data)
	if err != nil {
		return fmt.Errorf("解析 OpenAPI spec 失败: %w", err)
	}
	model, err := doc.BuildV3Model()
	if err != nil {
		return fmt.Errorf("构建 OpenAPI 模型失败: %w", err)
	}

	actions, err := collectActions(&model.Model)
	if err != nil {
		return err
	}
	if len(actions) == 0 {
		return fmt.Errorf("OpenAPI spec 中没有可生成的 action")
	}

	version := model.Model.Version
	if model.Model.Info != nil && model.Model.Info.Version != "" {
		version = model.Model.Info.Version
	}
	models := orderedmap.New[string, *highbase.SchemaProxy]()
	if model.Model.Components != nil && model.Model.Components.Schemas != nil {
		for name, schema := range model.Model.Components.Schemas.FromOldest() {
			models.Set(name, schema)
		}
	}
	for _, action := range actions {
		request, err := modelSchema(action.RequestSchema, action.Summary+" 请求参数。")
		if err != nil {
			return fmt.Errorf("生成 %s 请求模型失败: %w", action.Action, err)
		}
		response, err := modelSchema(action.ResponseSchema, action.Summary+" 响应数据。")
		if err != nil {
			return fmt.Errorf("生成 %s 响应模型失败: %w", action.Action, err)
		}
		if err := addModel(models, action.Name+"Request", request); err != nil {
			return err
		}
		if err := addModel(models, action.Name+"Response", response); err != nil {
			return err
		}
	}

	typeSource, modelTypes, err := generateModels(models, version)
	if err != nil {
		return fmt.Errorf("生成模型失败: %w", err)
	}
	for i := range actions {
		actions[i].RequestType = modelTypes[actions[i].Name+"Request"]
		actions[i].ResponseType = modelTypes[actions[i].Name+"Response"]
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}
	files := []struct {
		name    string
		content []byte
	}{
		{name: "actions_gen.go", content: generateActions(actions, version)},
		{name: "client_gen.go", content: generateClient(actions, version)},
		{name: "types_gen.go", content: typeSource},
	}
	for _, file := range files {
		formatted, err := format.Source(file.content)
		if err != nil {
			return fmt.Errorf("格式化 %s 失败: %w\n%s", file.name, err, file.content)
		}
		if err := os.WriteFile(filepath.Join(outDir, file.name), formatted, 0o644); err != nil {
			return fmt.Errorf("写入 %s 失败: %w", file.name, err)
		}
	}
	return nil
}

type actionSpec struct {
	Action         string
	Name           string
	Summary        string
	RequestType    string
	ResponseType   string
	RequestSchema  *highbase.SchemaProxy
	ResponseSchema *highbase.SchemaProxy
}

func collectActions(doc *highv3.Document) ([]actionSpec, error) {
	if doc == nil || doc.Paths == nil || doc.Paths.PathItems == nil {
		return nil, nil
	}
	paths := make([]string, 0, doc.Paths.PathItems.Len())
	for path := range doc.Paths.PathItems.FromOldest() {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	actions := make([]actionSpec, 0, len(paths))
	usedNames := map[string]int{}
	for _, path := range paths {
		item, _ := doc.Paths.PathItems.Get(path)
		if item == nil || item.Post == nil {
			continue
		}
		responseSchema, err := responseDataSchema(item.Post)
		if err != nil {
			return nil, fmt.Errorf("读取 %s 响应 schema 失败: %w", path, err)
		}
		action := strings.TrimPrefix(path, "/")
		actions = append(actions, actionSpec{
			Action:         action,
			Name:           uniqueName(ToExportedName(action), usedNames),
			Summary:        item.Post.Summary,
			RequestSchema:  requestSchema(item.Post),
			ResponseSchema: responseSchema,
		})
	}
	return actions, nil
}

func requestSchema(op *highv3.Operation) *highbase.SchemaProxy {
	if op == nil || op.RequestBody == nil || op.RequestBody.Content == nil {
		return nil
	}
	mediaType, ok := op.RequestBody.Content.Get("application/json")
	if !ok || mediaType == nil {
		return nil
	}
	return mediaType.Schema
}

func responseDataSchema(op *highv3.Operation) (*highbase.SchemaProxy, error) {
	if op == nil || op.Responses == nil {
		return nil, nil
	}
	response := op.Responses.FindResponseByCode(200)
	if response == nil {
		response = op.Responses.Default
	}
	if response == nil || response.Content == nil {
		return nil, nil
	}
	mediaType, ok := response.Content.Get("application/json")
	if !ok || mediaType == nil {
		return nil, nil
	}
	return findDataSchema(mediaType.Schema, make(map[*highbase.SchemaProxy]bool))
}

func findDataSchema(proxy *highbase.SchemaProxy, seen map[*highbase.SchemaProxy]bool) (*highbase.SchemaProxy, error) {
	if proxy == nil || seen[proxy] {
		return nil, nil
	}
	seen[proxy] = true
	schema, err := proxy.BuildSchema()
	if err != nil {
		return nil, err
	}
	if schema == nil {
		return nil, fmt.Errorf("schema 为空")
	}
	if schema.Properties != nil {
		if data, ok := schema.Properties.Get("data"); ok {
			return data, nil
		}
	}
	// allOf 后面的接口定义会覆盖前面的通用 envelope。
	for i := len(schema.AllOf) - 1; i >= 0; i-- {
		data, err := findDataSchema(schema.AllOf[i], seen)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return data, nil
		}
	}
	return nil, nil
}

func modelSchema(proxy *highbase.SchemaProxy, description string) (*highbase.SchemaProxy, error) {
	if proxy == nil {
		return highbase.CreateSchemaProxy(&highbase.Schema{
			Type:        []string{"object"},
			Description: description,
		}), nil
	}
	schema, err := proxy.BuildSchema()
	if err != nil {
		return nil, err
	}
	if schema == nil {
		return nil, fmt.Errorf("schema 为空")
	}
	clone := *schema
	if clone.Description == "" {
		clone.Description = description
	}
	return highbase.CreateSchemaProxy(&clone), nil
}

func addModel(models *orderedmap.Map[string, *highbase.SchemaProxy], name string, schema *highbase.SchemaProxy) error {
	if _, exists := models.Get(name); exists {
		return fmt.Errorf("生成模型名称冲突: %s", name)
	}
	models.Set(name, schema)
	return nil
}

func generateModels(schemas *orderedmap.Map[string, *highbase.SchemaProxy], version string) ([]byte, map[string]string, error) {
	generator := gogen.NewGenerator(
		gogen.WithPackageName("api"),
		gogen.WithGeneratedComment(true),
		gogen.WithHeaderComment(fmt.Sprintf("代码由 napcatgen 根据 NapCat OpenAPI %s 生成；请勿手动修改。", version)),
		gogen.WithNestedTypeNameDelimiter(""),
		gogen.WithTypeNameResolver(ToExportedName),
		gogen.WithEnumConstants(true),
	)
	file, err := generator.RenderSchemas(schemas)
	if err != nil {
		return nil, nil, err
	}
	if len(file.Types) != schemas.Len() {
		return nil, nil, fmt.Errorf("生成类型数量不匹配: schema=%d type=%d", schemas.Len(), len(file.Types))
	}
	types := make(map[string]string, len(file.Types))
	i := 0
	for name := range schemas.FromOldest() {
		generated := file.Types[i]
		typeName := generated.Name
		if generated.Kind == gogen.KindUnion {
			typeName += "Union"
		}
		types[name] = typeName
		i++
	}
	source, err := redirectJSONImport(file.Source)
	if err != nil {
		return nil, nil, err
	}
	return source, types, nil
}

func redirectJSONImport(source []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "types_gen.go", source, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("解析生成模型失败: %w", err)
	}
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, fmt.Errorf("解析生成模型 import 失败: %w", err)
		}
		if path == "encoding/json" {
			spec.Name = ast.NewIdent("json")
			spec.Path.Value = strconv.Quote("github.com/zjutjh/napcat-sdk/internal/jsonx")
		}
	}
	var out bytes.Buffer
	if err := format.Node(&out, fset, file); err != nil {
		return nil, fmt.Errorf("重写生成模型 JSON import 失败: %w", err)
	}
	return out.Bytes(), nil
}

func generateActions(actions []actionSpec, version string) []byte {
	var b bytes.Buffer
	writeHeader(&b, version)
	b.WriteString("package api\n\n")
	b.WriteString("// Action 是 NapCat API action 名称。\n")
	b.WriteString("type Action string\n\n")
	b.WriteString("const (\n")
	for _, action := range actions {
		writeComment(&b, "Action"+action.Name, action.Summary)
		fmt.Fprintf(&b, "\tAction%s Action = %q\n", action.Name, action.Action)
	}
	b.WriteString(")\n")
	return b.Bytes()
}

func generateClient(actions []actionSpec, version string) []byte {
	var b bytes.Buffer
	writeHeader(&b, version)
	b.WriteString("package api\n\n")
	b.WriteString("import \"context\"\n\n")
	for _, action := range actions {
		writeComment(&b, action.Name, action.Summary)
		fmt.Fprintf(&b, "func (c *Client) %s(ctx context.Context, req %s) (*%s, error) {\n", action.Name, action.RequestType, action.ResponseType)
		fmt.Fprintf(&b, "\tvar out %s\n", action.ResponseType)
		fmt.Fprintf(&b, "\tif err := c.caller.Call(ctx, string(Action%s), req, &out); err != nil {\n", action.Name)
		b.WriteString("\t\treturn nil, err\n")
		b.WriteString("\t}\n")
		b.WriteString("\treturn &out, nil\n")
		b.WriteString("}\n\n")
	}
	return b.Bytes()
}

func writeHeader(b *bytes.Buffer, version string) {
	fmt.Fprintf(b, "// 代码由 napcatgen 根据 NapCat OpenAPI %s 生成；请勿手动修改。\n\n", version)
}

func writeComment(b *bytes.Buffer, name string, text string) {
	text = sanitizeComment(text)
	if text == "" {
		fmt.Fprintf(b, "// %s 由 NapCat OpenAPI 生成。\n", name)
		return
	}
	fmt.Fprintf(b, "// %s %s\n", name, text)
}

func sanitizeComment(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	return text
}

func uniqueName(name string, used map[string]int) string {
	if name == "" {
		name = "Value"
	}
	used[name]++
	if used[name] == 1 {
		return name
	}
	return fmt.Sprintf("%s%d", name, used[name])
}

// ToExportedName 将 action 名称转换为导出的 Go 名称。
func ToExportedName(action string) string {
	prefix := ""
	if strings.HasPrefix(action, ".") {
		prefix = "Dot"
	}
	if strings.HasPrefix(action, "_") {
		prefix = "Underscore"
	}
	tokens := splitName(action)
	var b strings.Builder
	b.WriteString(prefix)
	for _, token := range tokens {
		b.WriteString(normalizeToken(token))
	}
	if b.Len() == 0 {
		return "Action"
	}
	return b.String()
}

var splitRegexp = regexp.MustCompile(`[^A-Za-z0-9]+`)

func splitName(name string) []string {
	raw := splitRegexp.Split(name, -1)
	out := make([]string, 0, len(raw))
	for _, token := range raw {
		if token != "" {
			out = append(out, token)
		}
	}
	return out
}

func normalizeToken(token string) string {
	upper := strings.ToUpper(token)
	switch upper {
	case "ID", "URL", "API", "JSON", "QQ", "UIN", "UID", "QID", "IP":
		return upper
	}
	runes := []rune(token)
	if len(runes) == 0 {
		return ""
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
