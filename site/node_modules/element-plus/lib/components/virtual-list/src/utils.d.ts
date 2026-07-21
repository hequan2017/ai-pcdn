import { Direction, RTLOffsetType } from "./types.js";
import { CSSProperties } from "vue";

//#region ../../packages/components/virtual-list/src/utils.d.ts
declare const getScrollDir: (prev: number, cur: number) => "forward" | "backward";
declare const isHorizontal: (dir: string) => dir is "ltr" | "rtl" | "horizontal";
declare const isRTL: (dir: Direction) => dir is "rtl";
declare function getRTLOffsetType(recalculate?: boolean): RTLOffsetType;
type RenderThumbStyleParams = {
  bar: {
    size: 'height' | 'width';
    axis: 'X' | 'Y';
  };
  size: string;
  move: number;
};
declare function renderThumbStyle({
  move,
  size,
  bar
}: RenderThumbStyleParams, layout: string): CSSProperties;
//#endregion
export { getRTLOffsetType, getScrollDir, isHorizontal, isRTL, renderThumbStyle };