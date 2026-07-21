import { Arrayable } from "../../../utils/typescript.js";
import { PopperEffect, PopperInstance } from "../../popper/src/popper.js";
import { Options, Placement } from "../../popper/index.js";
import { TooltipTriggerType } from "../../tooltip/src/trigger.js";
import { PopoverProps } from "./popover.js";
import * as _$vue from "vue";

//#region ../../packages/components/popover/src/popover.vue.d.ts
declare var __VLS_15: {}, __VLS_18: {
    hide: () => void;
  };
type __VLS_Slots = {} & {
  reference?: (props: typeof __VLS_15) => any;
} & {
  default?: (props: typeof __VLS_18) => any;
};
declare const __VLS_base: _$vue.DefineComponent<PopoverProps, {
  /** @description popper ref */popperRef: _$vue.ComputedRef<PopperInstance | undefined>; /** @description hide popover */
  hide: () => void;
}, {}, {}, {}, _$vue.ComponentOptionsMixin, _$vue.ComponentOptionsMixin, {
  "update:visible": (value: boolean) => void;
  "before-enter": () => void;
  "before-leave": () => void;
  "after-enter": () => void;
  "after-leave": () => void;
}, string, _$vue.PublicProps, Readonly<PopoverProps> & Readonly<{
  "onUpdate:visible"?: ((value: boolean) => any) | undefined;
  "onBefore-enter"?: (() => any) | undefined;
  "onBefore-leave"?: (() => any) | undefined;
  "onAfter-enter"?: (() => any) | undefined;
  "onAfter-leave"?: (() => any) | undefined;
}>, {
  tabindex: string | number;
  width: string | number;
  visible: boolean | null;
  content: string;
  teleported: boolean;
  popperStyle: string | false | _$vue.CSSProperties | _$vue.StyleValue[] | null;
  persistent: boolean;
  effect: PopperEffect;
  enterable: boolean;
  offset: number;
  placement: Placement;
  popperOptions: Partial<Options>;
  showArrow: boolean;
  showAfter: number;
  hideAfter: number;
  autoClose: number;
  trigger: Arrayable<TooltipTriggerType>;
  triggerKeys: string[];
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