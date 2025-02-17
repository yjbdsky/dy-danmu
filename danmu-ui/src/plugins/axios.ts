import axios from 'axios';
import router from '../router/router';
import { ElMessage } from 'element-plus';
import type { ApiResponse} from '../types/response';
import { computed } from 'vue';
import { useAuthStore } from '../stores/auth';
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json'
  }
});

api.interceptors.request.use(config => {
  const token = computed(() => useAuthStore().token);
  if (token) {
    config.headers.Authorization = `Bearer ${token.value}`
  }
  return config
})

api.interceptors.response.use(
  response => {
    const res = response.data as ApiResponse<any>;

    
    if (res.code === 200) {
      return response;
    }
    
    // 处理业务错误
    ElMessage.error(res.message || '操作失败');
    return Promise.reject(new Error(res.message || '操作失败'));
  },
  error => {
    // 处理 HTTP 错误
    if (error.response?.status === 401) {
      router.push('/login');
      ElMessage.error('请先登录');
    } else if (error.code === 'ECONNABORTED') {
      ElMessage.error('请求超时，请稍后重试');
    } else {
      ElMessage.error(error.message || '网络错误，请稍后重试');
    }
    return Promise.reject(error);
  }
)

export default api;
