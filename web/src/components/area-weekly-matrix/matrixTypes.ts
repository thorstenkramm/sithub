import type { MatrixCell, MatrixItem } from '../../api/itemGroupMatrix';

export interface MatrixCellClickEvent {
  type: 'book' | 'cancel';
  el: HTMLElement;
  item: MatrixItem;
  cell: MatrixCell;
}

export type MatrixCellClickHandler = (event: MatrixCellClickEvent) => void;
