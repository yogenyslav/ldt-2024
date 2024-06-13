// Import any necessary libraries or modules

// Define the login form HTML
const loginForm = `
<form id="login-form">
  <label for="username">Username:</label>
  <input type="text" id="username" name="username" required><br><br>
  
  <label for="password">Password:</label>
  <input type="password" id="password" name="password" required><br><br>
  
  <button type="submit">Login</button>
</form>
`;

// Add the login form to the DOM
document.body.innerHTML = loginForm;

// Handle form submission
document.getElementById('login-form').addEventListener('submit', async (event) => {
  event.preventDefault();

  // Get the entered username and password
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;
  const urlParams = new URLSearchParams(window.location.search);
  const tg_id = Number(urlParams.get('tg_id'));

  // Send a request to the API to validate the credentials
  try {
    const response = await fetch('http://localhost:9998/api/v1/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username, password })
    });

    if (response.ok) {
      // Credentials are valid, do something
      console.log('Login successful!');
      const data = await response.json();
      console.log('Data:', data);

      try {
      const resp = await fetch('http://localhost:11000/bot/auth', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ tg_id, token: data.token, roles: data.roles })
      })
    } catch (error) {
      console.error('An error occurred:', error);
    }
    } else {
      // Credentials are invalid, do something else
      console.log('Invalid credentials!');
    }
  } catch (error) {
    console.error('An error occurred:', error);
  }

  
});