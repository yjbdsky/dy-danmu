export const enum Role {
  Admin = 'admin',
  User = 'guest',
}

//Post /api/auth/login
export interface LoginRequest {
  email: string;
  password: string;
}

//Post /api/auth/register admin
export interface RegisterRequest {
  name: string;
  password: string;
  email: string;
  role: Role;
}

//PUT /api/auth/self
export interface UpdateUserRequest {
  email: string;
  name: string;
  password?: string;
}

export interface Auth {
  id: string;
  name: string;
  email: string;
  role: Role;
}

export interface TokenResponse {
  token: string;
}


