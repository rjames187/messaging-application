import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';

const API_URL = import.meta.env.VITE_API_URL;

const View = () => {
  const [user, setUser] = useState({
    firstName: "",
    lastName: ""
  });
  const navigate = useNavigate();

  useEffect(() => {
    const authHeader = localStorage.getItem('authorization') ?? '';
    if (!authHeader.startsWith("Bearer ")) {
      navigate('/');
    }

    (async () => {
      try {
        const response = await fetch(`${API_URL}/v1/users/me`, {
        headers: {
          'Authorization': authHeader
        }
        });

        if (!response.ok) {
          throw new Error('Failed to fetch user data');
        }

        const data = await response.json();
        setUser({
          firstName: data.firstName,
          lastName: data.lastName
        });
      } catch (error) {
        console.error('Error fetching user data:', error);
        alert(error.message);
      }
    })();
  }, [])

  const handleSignOut = async () => {
    console.log('Signing out...');

    try {
      const response = await fetch(`${API_URL}/v1/sessions/mine`, {
        method: 'DELETE',
        headers: {
          'Authorization': localStorage.getItem('authorization')
        }
      });

      if (!response.ok) {
        throw new Error('Sign out failed');
      }

      console.log('Sign out successful');
    } catch (error) {
      console.error('Error during sign out:', error);
      alert(error.message);
      return;
    }

    localStorage.removeItem('authorization');
    navigate('/');
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
          <Link to="/update" className="update-profile-link">
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

export default View;