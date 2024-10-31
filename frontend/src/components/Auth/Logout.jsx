export default function Logout() {
  const handleLogout = () => {
    // Add logout functionality here, possibly via an API call
  };

  return (
    <button
      onClick={handleLogout}
      className="mt-4 w-full rounded-md bg-red-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-red-500 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
    >
      Logout
    </button>
  );
}
