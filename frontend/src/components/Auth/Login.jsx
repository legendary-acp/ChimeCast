import { useState } from 'react';
import { login } from "../../services/authService"; // Make sure to import your login function

export default function Login() {
  const [formData, setFormData] = useState({
    username: '', // Change to username
    password: '', // Change to password
  });

  const [error, setError] = useState(null); // To handle error messages
  const [success, setSuccess] = useState(null); // To handle success messages

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault(); // Prevent default form submission

    setError(null); // Reset error message
    setSuccess(null); // Reset success message

    try {
      const response = await login(formData.username, formData.password); // Call your login function
      setSuccess(response.message || "Login successful!"); // Handle success message
      // Optionally, redirect the user or clear the form
    } catch (error) {
      // Set error to a string
      setError(error.message || "An unknown error occurred"); // Ensure error is a string
    }
  };

  return (
    <div>
      <div className="flex flex-col items-center">
        <img alt="Chime Cast" src="/image.png" className="h-20 w-auto" />
        <h2 className="mt-6 text-center text-2xl font-bold tracking-tight text-gray-900">
          Sign in to your account
        </h2>
      </div>

      <form onSubmit={handleSubmit} className="mt-8 space-y-6">
        {error && <p className="text-red-500">{error}</p>} {/* Display error message */}
        {success && <p className="text-green-500">{success}</p>} {/* Display success message */}

        <div>
          <label htmlFor="username" className="block text-sm font-medium text-gray-900">
            Username
          </label>
          <input
            id="username"
            name="username" // Ensure this matches formData key
            type="text" // Change to text
            required
            autoComplete="username"
            value={formData.username} // Ensure this matches formData key
            onChange={handleChange}
            className="mt-2 block w-full rounded-md border-gray-300 px-3 py-2 text-gray-900 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
        </div>

        <div>
          <label htmlFor="password" className="block text-sm font-medium text-gray-900">
            Password
          </label>
          <input
            id="password"
            name="password" // Ensure this matches formData key
            type="password"
            required
            autoComplete="current-password"
            value={formData.password} // Ensure this matches formData key
            onChange={handleChange}
            className="mt-2 block w-full rounded-md border-gray-300 px-3 py-2 text-gray-900 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
        </div>

        <button
          type="submit"
          className="mt-4 w-full rounded-md bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
        >
          Sign in
        </button>
      </form>
    </div>
  );
}
