import api from '../plugins/axios';
import type { LoginRequest, RegisterRequest, UpdateUserRequest, TokenResponse, Auth } from '../types/models/auth';
import type {ApiResponse} from "../types/response.ts";

export function login(data: LoginRequest) {
  return api.post<ApiResponse<TokenResponse>>('/api/auth/login', data);
}

export function register(data: RegisterRequest) {
  return api.post<ApiResponse<null>>('/api/auth/register', data);
}

export function resetPassword(id: string) {
  return api.post<ApiResponse<null>>(`/api/auth/reset-password/${id}`);
}

export function deleteUser(id: string) {
  return api.delete<ApiResponse<null>>(`/api/auth/${id}`);
}

export function updateUser(data: UpdateUserRequest) {
  return api.put<ApiResponse<null>>(`/api/auth/self`, data);
}

export function listUsers() {
  return api.get<ApiResponse<Auth[]>>('/api/auth/list');
}

export function getAuth() {
  return api.get<ApiResponse<Auth>>('/api/auth/self');
}


