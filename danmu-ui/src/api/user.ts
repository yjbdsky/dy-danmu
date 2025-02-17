import type { User } from "../types/models/user";

import api from '../plugins/axios';
import type { UserPageRequest, UserSearchRequest } from '../types/models/user';
import type { ApiResponse, Page } from '../types/response';

export function listAllUser(data: UserPageRequest) {
    return api.get<ApiResponse<Page<User>>>('/api/user', { params: data });
}

export function searchUser(data: UserSearchRequest) {
    return api.get<ApiResponse<Page<User>>>('/api/user/search', { params: data });
}


