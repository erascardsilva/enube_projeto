import authService from "../services/authService";

const userApi = {
  login: async (username, password) => {
    try {
      const response = await fetch("http://localhost:8080/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
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
      const response = await fetch("http://localhost:8080/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, email, password }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error("Erro no registro:", error);
      throw error;
    }
  },

  // Exemplo de uma chamada autenticada
  getProtectedData: async () => {
    try {
      const token = authService.getToken();
      const headers = {
        "Content-Type": "application/json",
      };
      if (token) {
        headers.Authorization = `Bearer ${token}`;
      }

      const response = await fetch("http://localhost:8080/some-protected-endpoint", {
        method: "GET",
        headers,
      });

      if (response.status === 401) {
        authService.logout();
        window.location.href = "/authentication/sign-in"; // Redireciona para o login
        return null; // Retorna nulo para indicar que a requisição falhou devido a autenticação
      }

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error("Erro ao buscar dados protegidos:", error);
      throw error;
    }
  },
};

export default userApi; 