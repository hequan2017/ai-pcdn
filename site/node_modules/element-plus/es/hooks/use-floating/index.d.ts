import { Ref, ToRefs } from "vue";
import * as _$_floating_ui_dom0 from "@floating-ui/dom";
import { Middleware, Placement, SideObject, Strategy, VirtualElement } from "@floating-ui/dom";

//#region ../../packages/hooks/use-floating/index.d.ts
declare const useFloatingProps: {};
type UseFloatingProps = ToRefs<{
  middleware: Array<Middleware>;
  placement: Placement;
  strategy: Strategy;
}>;
declare const getPositionDataWithUnit: <T extends Record<string, number>>(record: T | undefined, key: keyof T) => string;
declare const useFloating: ({
  middleware,
  placement,
  strategy
}: UseFloatingProps) => {
  update: () => Promise<void>;
  referenceRef: Ref<HTMLElement | VirtualElement | undefined, HTMLElement | VirtualElement | undefined>;
  contentRef: Ref<HTMLElement | undefined, HTMLElement | undefined>;
  x: Ref<number | undefined, number | undefined>;
  y: Ref<number | undefined, number | undefined>;
  placement: Ref<Placement, Placement>;
  strategy: Ref<Strategy, Strategy>;
  middlewareData: Ref<{
    [x: string]: any;
    arrow?: {
      y?: number | undefined;
      x?: number | undefined;
      centerOffset: number;
      alignmentOffset?: number | undefined;
    } | undefined;
    autoPlacement?: {
      index?: number | undefined;
      overflows: {
        placement: Placement;
        overflows: Array<number>;
      }[];
    } | undefined;
    flip?: {
      index?: number | undefined;
      overflows: {
        placement: Placement;
        overflows: Array<number>;
      }[];
    } | undefined;
    hide?: {
      referenceHidden?: boolean | undefined;
      escaped?: boolean | undefined;
      referenceHiddenOffsets?: {
        left: number;
        right: number;
        bottom: number;
        top: number;
      } | undefined;
      escapedOffsets?: {
        left: number;
        right: number;
        bottom: number;
        top: number;
      } | undefined;
    } | undefined;
    offset?: {
      y: number;
      x: number;
      placement: Placement;
    } | undefined;
    shift?: {
      y: number;
      x: number;
      enabled: {
        y: boolean;
        x: boolean;
      };
    } | undefined;
  }, _$_floating_ui_dom0.MiddlewareData | {
    [x: string]: any;
    arrow?: {
      y?: number | undefined;
      x?: number | undefined;
      centerOffset: number;
      alignmentOffset?: number | undefined;
    } | undefined;
    autoPlacement?: {
      index?: number | undefined;
      overflows: {
        placement: Placement;
        overflows: Array<number>;
      }[];
    } | undefined;
    flip?: {
      index?: number | undefined;
      overflows: {
        placement: Placement;
        overflows: Array<number>;
      }[];
    } | undefined;
    hide?: {
      referenceHidden?: boolean | undefined;
      escaped?: boolean | undefined;
      referenceHiddenOffsets?: {
        left: number;
        right: number;
        bottom: number;
        top: number;
      } | undefined;
      escapedOffsets?: {
        left: number;
        right: number;
        bottom: number;
        top: number;
      } | undefined;
    } | undefined;
    offset?: {
      y: number;
      x: number;
      placement: Placement;
    } | undefined;
    shift?: {
      y: number;
      x: number;
      enabled: {
        y: boolean;
        x: boolean;
      };
    } | undefined;
  }>;
};
type ArrowMiddlewareProps = {
  arrowRef: Ref<HTMLElement | null | undefined>;
  padding?: number | SideObject;
};
declare const arrowMiddleware: ({
  arrowRef,
  padding
}: ArrowMiddlewareProps) => Middleware;
//#endregion
export { ArrowMiddlewareProps, UseFloatingProps, arrowMiddleware, getPositionDataWithUnit, useFloating, useFloatingProps };