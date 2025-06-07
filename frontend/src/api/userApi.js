import axios from 'axios';

const baseUrl = 'http://localhost:8080';

const userApi = {
  login: async (username, password) => {
    try {
      const response = await axios.post(`${baseUrl}/auth/login`, { username, password });
      const data = response.data;
      if (data.token) {
        localStorage.setItem("user", JSON.stringify(data));
      }
      return data;
    } catch (error) {
      console.error("Erro no login:", error);
      throw error;
    }
  },

  register: async (username, email, password) => {
    try {
      const response = await axios.post(`${baseUrl}/auth/register`, { username, email, password });
      return response.data;
    } catch (error) {
      console.error("Erro no registro:", error);
      throw error;
    }
  },

  getProtectedData: async () => {
    try {
      const token = localStorage.getItem("user") ? JSON.parse(localStorage.getItem("user")).token : null;
      if (!token) {
        throw new Error("No authentication token found");
      }
      const response = await axios.get(`${baseUrl}/some-protected-endpoint`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      return response.data;
    } catch (error) {
      console.error("Erro ao buscar dados protegidos:", error);
      throw error;
    }
  },

  fetchUsers: async (page = 1, limit = 10) => {
    try {
      const token = localStorage.getItem("user") ? JSON.parse(localStorage.getItem("user")).token : null;
      if (!token) {
        throw new Error("No authentication token found");
      }
      const response = await axios.get(`${baseUrl}/api/users?page=${page}&limit=${limit}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      return response.data;
    } catch (error) {
      throw error;
    }
  },
};

export default userApi; 