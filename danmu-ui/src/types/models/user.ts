export interface User {
  id: number;
  user_id: number;
  user_name: string;
  display_id: string;
}

export interface UserPageRequest {
  page: number;
  page_size: number;
}

export interface UserSearchRequest {
  keyword: string;
  page: number;
  page_size: number;
}
