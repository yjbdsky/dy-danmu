import { defineStore } from 'pinia';
import { ref } from 'vue';
import { login as loginApi, getAuth as getAuthApi } from '../api/auth';
import type { Auth, LoginRequest } from '../types/models/auth';
import router from "../router/router.ts";

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'));
  const auth = ref<Auth | null>(JSON.parse(localStorage.getItem('auth') || 'null'));
  const isAuthenticated = ref(!!token.value);

  async function login(loginData: LoginRequest) {
    try {
      const tokenResponse = (await loginApi(loginData)).data;
      token.value = tokenResponse.data.token;
      const authResponse = (await getAuthApi()).data;
      auth.value = authResponse.data;
      console.log(auth.value.role)
      isAuthenticated.value = true;

      localStorage.setItem('token', token.value);
      localStorage.setItem('auth', JSON.stringify(auth.value));
    } catch (error) {
      token.value = null;
      auth.value = null;
      isAuthenticated.value = false;
      
      localStorage.removeItem('token');
      localStorage.removeItem('auth');
      
      throw error;
    }
  }

  function logout() {
    token.value = null;
    auth.value = null;
    isAuthenticated.value = false;
    
    localStorage.removeItem('token');
    localStorage.removeItem('auth');
    router.push("/login")
  }

  return {
    token,
    auth,
    isAuthenticated,
    login,
    logout,
  };
});
