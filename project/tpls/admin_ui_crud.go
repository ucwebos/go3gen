package tpls

const AdminUICrud = `import { CreateCrudOptionsProps, CreateCrudOptionsRet, dict } from "@fast-crud/fast-crud";
import { addRequest, delRequest, editRequest, pageRequest } from "./api";

export default function ({ crudExpose, context }: CreateCrudOptionsProps): CreateCrudOptionsRet {
  return {
    crudOptions: {
      // 自定义crudOptions配置
      request: {
        pageRequest,
        addRequest,
        editRequest,
        delRequest
      },
      // 字段配置
      columns: {
{{- range .Fields}}
        {{.Name}}: {
          title: "{{.Title}}",
          type: "text",
          search: { show: true },
          column: {
            show: true,
            resizable: true,
            width: 100
          }
        },
{{- end}}
      }
    }
  };
}
`
