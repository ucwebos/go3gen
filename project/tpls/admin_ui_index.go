package tpls

const AdminUIIndex = `
<template>
    <fs-page>
        <fs-crud ref="crudRef" v-bind="crudBinding" />
    </fs-page>
</template>

<script lang="ts">
    import { defineComponent, onMounted } from "vue";
    import { useFs ,OnExposeContext } from "@fast-crud/fast-crud";
    import createCrudOptions from "./crud";

    //此处为组件定义
    export default defineComponent({
        name: "{{.Name}}",
        setup(props:any,ctx:any) {
            const context: any = {props,ctx}; // 自定义变量, 将会传递给createCrudOptions, 比如直接把props,和ctx直接传过去使用
            function onExpose(e:OnExposeContext){} //将在createOptions之前触发，可以获取到crudExpose,和context
            const { crudRef, crudBinding, crudExpose } = useFs({ createCrudOptions, onExpose, context});
            // 页面打开后获取列表数据
            onMounted(() => {
                crudExpose.doRefresh();
            });
            return {
                crudBinding,
                crudRef
            };
        }
    });
</script>
`
