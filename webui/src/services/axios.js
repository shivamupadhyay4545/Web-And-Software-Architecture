import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

export const setAuthToken = (token) => {
	if (token) {
	  instance.defaults.headers.common["Authorization"] = token;
	  // Store the token in local storage or another secure storage mechanism if needed
	  localStorage.setItem("authToken", token);
	} else {
	  delete instance.defaults.headers.common["Authorization"];
	  // Remove the token from storage if it's cleared
	  localStorage.removeItem("authToken");
	}
  };
  
  // Check if there's a stored token on initialization
  const storedToken = localStorage.getItem("authToken");
  if (storedToken) {
	setAuthToken(storedToken);
  }
  

export default instance;
