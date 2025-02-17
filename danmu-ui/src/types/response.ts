export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

export interface Page<T> {
  total: number;
  list: T[];
}



