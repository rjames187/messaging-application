import { Link } from 'react-router-dom';

const Update = () => {
  // Mock user data - replace with actual user data from context/props/API
  const user = {
    firstName: 'John',
    lastName: 'Doe'
  };

  const handleSignOut = () => {
    // Add signout logic here (clear tokens, redirect, etc.)
    console.log('Signing out...');
    // Example: localStorage.removeItem('token');
    // Example: navigate('/login');
  };

  return (
    <div className="view-container">
      <div className="user-info">
        <h1>Welcome, {user.firstName} {user.lastName}</h1>
        
        <div className="user-details">
          <p><strong>First Name:</strong> {user.firstName}</p>
          <p><strong>Last Name:</strong> {user.lastName}</p>
        </div>
        
        <div className="actions">
          <Link to="/profile/edit" className="update-profile-link">
            Update Profile
          </Link>
          
          <button 
            onClick={handleSignOut}
            className="signout-button"
            type="button"
          >
            Sign Out
          </button>
        </div>
      </div>
    </div>
  );
};

export default Update;