import { CheckboxProps, CheckboxValueType } from "./checkbox.js";
import * as _$vue from "vue";

//#region ../../packages/components/checkbox/src/checkbox.vue.d.ts
declare var __VLS_10: {};
type __VLS_Slots = {} & {
  default?: (props: typeof __VLS_10) => any;
};
declare const __VLS_base: _$vue.DefineComponent<CheckboxProps, {}, {}, {}, {}, _$vue.ComponentOptionsMixin, _$vue.ComponentOptionsMixin, {
  "update:modelValue": (val: CheckboxValueType) => void;
  change: (val: CheckboxValueType) => void;
}, string, _$vue.PublicProps, Readonly<CheckboxProps> & Readonly<{
  "onUpdate:modelValue"?: ((val: CheckboxValueType) => any) | undefined;
  onChange?: ((val: CheckboxValueType) => any) | undefined;
}>, {
  modelValue: number | string | boolean;
  label: string | boolean | number | object;
  disabled: boolean;
  value: string | boolean | number | object;
  name: string;
  trueValue: string | number;
  falseValue: string | number;
  trueLabel: string | number;
  falseLabel: string | number;
  id: string;
  validateEvent: boolean;
}, {}, {}, {}, string, _$vue.ComponentProvideOptions, false, {}, any>;
declare const __VLS_export: __VLS_WithSlots<typeof __VLS_base, __VLS_Slots>;
declare const _default: typeof __VLS_export;
type __VLS_WithSlots<T, S> = T & {
  new (): {
    $slots: S;
  };
};
//#endregion
export { _default as default };