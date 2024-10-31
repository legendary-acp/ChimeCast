import { useState } from 'react';
import Login from './Login';
import Register from './Register';

export default function ToggleAuth() {
  const [isRegister, setIsRegister] = useState(true);

  const toggleScreen = () => setIsRegister(!isRegister);

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-8">
      <div className="w-full max-w-md space-y-8 rounded-lg bg-white p-8 shadow-lg">
        {isRegister ? <Register /> : <Login />}
        <div className="flex justify-center">
          <button
            onClick={toggleScreen}
            className="text-sm font-semibold text-indigo-600 hover:text-indigo-500"
          >
            {isRegister ? 'Already have an account? Sign in' : "Don't have an account? Register"}
          </button>
        </div>
      </div>
    </div>
  );
}
