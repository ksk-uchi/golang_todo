export interface Todo {
  id: number;
  title: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Pagination {
  total_pages: number;
  current_page: number;
  has_next: boolean;
  has_prev: boolean;
  limit: number;
}

export interface ListTodoResponse {
  data: Todo[];
  pagination: Pagination;
}
