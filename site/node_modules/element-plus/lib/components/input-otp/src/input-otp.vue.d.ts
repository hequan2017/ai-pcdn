import { InputOtpProps } from "./input-otp.js";
import * as _$vue from "vue";

//#region ../../packages/components/input-otp/src/input-otp.vue.d.ts
declare var __VLS_1: {
  index: number;
};
type __VLS_Slots = {} & {
  separator?: (props: typeof __VLS_1) => any;
};
declare const __VLS_base: _$vue.DefineComponent<InputOtpProps, {
  /**
   * @description HTML input elements array
   */
  inputRefs: _$vue.Ref<(HTMLInputElement | undefined)[], (HTMLInputElement | undefined)[]>;
  /**
   * @description Focus an OTP input field
   */
  focus: (index?: number) => void;
  /**
   * @description Blur the focused OTP input field
   */
  blur: () => void;
}, {}, {}, {}, _$vue.ComponentOptionsMixin, _$vue.ComponentOptionsMixin, {
  "update:modelValue": (value: string) => void;
  finish: (value: string) => void;
  change: (value: string) => void;
  blur: (eve: FocusEvent) => void;
  focus: (eve: FocusEvent) => void;
}, string, _$vue.PublicProps, Readonly<InputOtpProps> & Readonly<{
  "onUpdate:modelValue"?: ((value: string) => any) | undefined;
  onFinish?: ((value: string) => any) | undefined;
  onChange?: ((value: string) => any) | undefined;
  onBlur?: ((eve: FocusEvent) => any) | undefined;
  onFocus?: ((eve: FocusEvent) => any) | undefined;
}>, {
  type: "outlined" | "filled" | "underlined";
  length: number;
  disabled: boolean;
  validateEvent: boolean;
  validator: (char: string, index: number) => boolean;
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