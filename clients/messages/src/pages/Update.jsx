import { useState } from 'react';
import { Link } from 'react-router-dom';

const API_URL = import.meta.env.VITE_API_URL;

const Update = () => {
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: ''
  });
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState('');

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (formData.firstName === '' && formData.lastName === '') {
      alert('Please enter at least one field to update.');
      return;
    }

    setIsLoading(true);
    setMessage('');

    const updates = {
      firstName: formData.firstName !== "" ? formData.firstName : undefined,
      lastName: formData.lastName !== "" ? formData.lastName : undefined
    }

    try {
      const response = await fetch(`${API_URL}/v1/users/me`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': localStorage.getItem('authorization')
        },
        body: JSON.stringify(updates)
      });

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      setMessage('Profile updated successfully!');
      setFormData({ firstName: '', lastName: '' });
    } catch (error) {
      setMessage('Error updating profile. Please try again.');
      console.error('There was a problem with the fetch operation:', error);
      alert(error.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="update-profile">
      <div className="container">
        <h1>Update Profile</h1>
        
        <form onSubmit={handleSubmit} className="profile-form">
          <div className="form-group">
            <label htmlFor="firstName">First Name:</label>
            <input
              type="text"
              id="firstName"
              name="firstName"
              value={formData.firstName}
              onChange={handleInputChange}
              disabled={isLoading}
            />
          </div>

          <div className="form-group">
            <label htmlFor="lastName">Last Name:</label>
            <input
              type="text"
              id="lastName"
              name="lastName"
              value={formData.lastName}
              onChange={handleInputChange}
              disabled={isLoading}
            />
          </div>

          <div className="form-actions">
            <button 
              type="submit" 
              disabled={isLoading}
              className="btn btn-primary"
            >
              {isLoading ? 'Updating...' : 'Update Profile'}
            </button>
            
            <Link to="/view" className="btn btn-secondary">
              View Profile
            </Link>
          </div>
        </form>

        {message && (
          <div className={`message ${message.includes('Error') ? 'error' : 'success'}`}>
            {message}
          </div>
        )}
      </div>
    </div>
  );
};

export default Update;