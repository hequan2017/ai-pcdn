import { PopperEffect } from "./popper.js";
import { PopperContentProps } from "./content.js";
import { Options as Options$1, Placement as Placement$1 } from "../index.js";
import * as _$vue from "vue";
import * as _$_popperjs_core0 from "@popperjs/core";

//#region ../../packages/components/popper/src/content.vue.d.ts
declare var __VLS_13: {};
type __VLS_Slots = {} & {
  default?: (props: typeof __VLS_13) => any;
};
declare const __VLS_base: _$vue.DefineComponent<PopperContentProps, {
  /**
   * @description popper content element
   */
  popperContentRef: _$vue.Ref<HTMLElement | undefined, HTMLElement | undefined>;
  /**
   * @description popperjs instance
   */
  popperInstanceRef: _$vue.ComputedRef<_$_popperjs_core0.Instance | undefined>;
  /**
   * @description method for updating popper
   */
  updatePopper: (shouldUpdateZIndex?: boolean) => void;
  /**
   * @description content style
   */
  contentStyle: _$vue.ComputedRef<_$vue.StyleValue[]>;
}, {}, {}, {}, _$vue.ComponentOptionsMixin, _$vue.ComponentOptionsMixin, {
  close: () => void;
  blur: () => void;
  focus: () => void;
  mouseenter: (evt: MouseEvent) => void;
  mouseleave: (evt: MouseEvent) => void;
}, string, _$vue.PublicProps, Readonly<PopperContentProps> & Readonly<{
  onClose?: (() => any) | undefined;
  onBlur?: (() => any) | undefined;
  onFocus?: (() => any) | undefined;
  onMouseenter?: ((evt: MouseEvent) => any) | undefined;
  onMouseleave?: ((evt: MouseEvent) => any) | undefined;
}>, {
  style: string | false | _$vue.CSSProperties | _$vue.StyleValue[] | null;
  visible: boolean;
  popperStyle: string | false | _$vue.CSSProperties | _$vue.StyleValue[] | null;
  effect: PopperEffect;
  enterable: boolean;
  pure: boolean;
  focusOnShow: boolean;
  trapping: boolean;
  stopPopperMouseEvent: boolean;
  virtualTriggering: boolean;
  loop: boolean;
  boundariesPadding: number;
  gpuAcceleration: boolean;
  offset: number;
  placement: Placement$1;
  popperOptions: Partial<Options$1>;
  strategy: "absolute" | "fixed";
  arrowOffset: number;
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