package tpls

import (
	"bytes"
	"text/template"
)

const AdminUIApiTpl = `import { AddReq, DelReq, EditReq, UserPageQuery, UserPageRes } from "@fast-crud/fast-crud";
import { request } from "/@/api/service";

const records = [{ id: 1, name: "Hello World", type: 1 }]
export const pageRequest = async (query: UserPageQuery): Promise<UserPageRes> => {
  return await request({
    url: "{{.ApiBaseURL}}/{{.NameVal}}-list",
    method: "post",
    data: query
  });
};
export const editRequest = async ({ form, row }: EditReq) => {
  if (form.id == null) {
  	form.id = row.id;
  }
  return await request({
    url: "{{.ApiBaseURL}}/{{.NameVal}}-edit",
    method: "post",
    data: form
  });
};
export const delRequest = async ({ row }: DelReq) => {
  return await request({
    url: "{{.ApiBaseURL}}/{{.NameVal}}-delete",
    method: "post",
    data: {id: row.id}
  });
};
export const addRequest = async ({ form }: AddReq) => {
  return await request({
    url: "{{.ApiBaseURL}}/{{.NameVal}}-add",
    method: "post",
    data: form
  });
};`

func (s *CrudItem) Execute(t string) ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("AdminUI").Parse(t)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
