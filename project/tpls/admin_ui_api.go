package tpls

import (
	"bytes"
	"text/template"
)

const AdminUIApiTpl = `import { AddReq, DelReq, EditReq, UserPageQuery, UserPageRes } from "@fast-crud/fast-crud";
import _ from "lodash-es";

const records = [{ id: 1, name: "Hello World", type: 1 }]
export const pageRequest = async (query: UserPageQuery): Promise<UserPageRes> => {
    return {
        records:_.cloneDeep(records), 
        offset: 0,
        limit: 20,
        total: records.length
    };
};
export const editRequest = async ({ form, row }: EditReq) => {
    const target = _.find(records, (item) => {
        return row.id === item.id;
    });
    _.merge(target, form);
    return target;
};
export const delRequest = async ({ row }: DelReq) => {
    _.remove(records, (item) => {
        return item.id === row.id;
    });
};
export const addRequest = async ({ form }: AddReq) => {
    const maxRecord = _.maxBy(records, (item) => {
        return item.id;
    });
    form.id = (maxRecord?.id || 0) + 1;
    records.push(form);
    return form;
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
