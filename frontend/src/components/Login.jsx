import React, { useState, useEffect } from 'react';
import { callApi, updateJwtToken } from '../utils/api';

const Login = ({ activeModel, setApiResponse }) => {
    const [activeLoginSection, setActiveLoginSection] = useState('authenticate');
    const [jwtToken, setJwtToken] = useState(localStorage.getItem("jwtToken") || "");

    useEffect(() => {
        updateJwtDisplay();
    }, [jwtToken]);

    const updateJwtDisplay = () => {
        const currentJwtDisplay = document.getElementById("currentJwt");
        if (currentJwtDisplay) {
            currentJwtDisplay.textContent = jwtToken || "No JWT token found.";
        }
    };

    const clearJwt = () => {
        updateJwtToken("");
        setJwtToken("");
        alert("JWT token cleared.");
    };

    const showLoginSection = (sectionName) => {
        setActiveLoginSection(sectionName);
        setApiResponse(''); // Clear API response when switching sections
    };

    const handleLoginFormSubmit = async (e) => {
        e.preventDefault();
        const email = e.target.loginEmail.value;
        const password = e.target.loginPassword.value;

        const data = await callApi("POST", `/login`, { email, password }, setApiResponse);
        if (data && data.token) {
            updateJwtToken(data.token);
            setJwtToken(data.token);
            alert("Login successful! JWT token stored.");
        } else {
            alert("Login failed.");
        }
    };

    const handleRegisterCustomerFormSubmit = async (e) => {
        e.preventDefault();
        const name = e.target.registerCustomerName.value;
        const email = e.target.registerCustomerEmail.value;
        const password = e.target.registerCustomerPassword.value;

        const data = await callApi("POST", `/register`, { name, email, password }, setApiResponse);
        if (data) {
            alert("Customer registration successful!");
        } else {
            alert("Customer registration failed.");
        }
    };

    const handleForgotPasswordFormSubmit = async (e) => {
        e.preventDefault();
        const email = e.target.forgotPasswordEmail.value;
        await callApi("POST", `/forgot-password`, { email }, setApiResponse);
    };

    const handleResetPasswordFormSubmit = async (e) => {
        e.preventDefault();
        const token = e.target.resetToken.value;
        const new_password = e.target.newPassword.value;
        await callApi("POST", `/reset-password`, { token, new_password }, setApiResponse);
    };

    if (activeModel !== 'login') {
        return null;
    }

    return (
        <div id="login" className="model-section active">
            <h2>Login / Register</h2>

            <div className="menu">
                <button onClick={() => showLoginSection('register-customer')} className={activeLoginSection === 'register-customer' ? 'active' : ''}>
                    Register
                </button>
                <button onClick={() => showLoginSection('authenticate')} className={activeLoginSection === 'authenticate' ? 'active' : ''}>
                    Login
                </button>
                <button onClick={() => showLoginSection('clear-jwt')}>Logout</button>
                <button onClick={() => showLoginSection('forgot-password')}>
                    Forgot Password
                </button>
            </div>

            {activeLoginSection === 'authenticate' && (
                <div id="login-authenticate" className="login-sub-section active">
                    <h3>Authenticate (POST /login)</h3>
                    <form id="loginForm" onSubmit={handleLoginFormSubmit}>
                        <div className="form-group-item">
                            <label htmlFor="loginEmail">Email:</label>
                            <input type="text" id="loginEmail" name="loginEmail" required />
                        </div>
                        <div className="form-group-item">
                            <label htmlFor="loginPassword">Password:</label>
                            <input type="password" id="loginPassword" name="loginPassword" required />
                        </div>
                        <button type="submit">Login</button>
                    </form>
                </div>
            )}

            {activeLoginSection === 'register-customer' && (
                <div id="login-register-customer" className="login-sub-section">
                    <h3>Register New Customer (POST /register)</h3>
                    <form id="registerCustomerForm" onSubmit={handleRegisterCustomerFormSubmit}>
                        <div className="form-group-item">
                            <label htmlFor="registerCustomerName">Name:</label>
                            <input type="text" id="registerCustomerName" name="registerCustomerName" required />
                        </div>
                        <div className="form-group-item">
                            <label htmlFor="registerCustomerEmail">Email:</label>
                            <input type="text" id="registerCustomerEmail" name="registerCustomerEmail" required />
                        </div>
                        <div className="form-group-item">
                            <label htmlFor="registerCustomerPassword">Password:</label>
                            <input type="password" id="registerCustomerPassword" name="registerCustomerPassword" required />
                        </div>
                        <button type="submit">Register Customer</button>
                    </form>
                </div>
            )}

            {activeLoginSection === 'clear-jwt' && (
                <div id="login-clear-jwt" className="login-sub-section">
                    <h3>Current JWT:</h3>
                    <pre id="currentJwt"></pre>
                    <button onClick={clearJwt}>Clear JWT</button>
                </div>
            )}

            {activeLoginSection === 'forgot-password' && (
                <div id="login-forgot-password" className="login-sub-section">
                    <h3>Forgot Password (POST /forgot-password)</h3>
                    <form id="forgotPasswordForm" onSubmit={handleForgotPasswordFormSubmit}>
                        <div className="form-group-item">
                            <label htmlFor="forgotPasswordEmail">Email:</label>
                            <input type="text" id="forgotPasswordEmail" name="forgotPasswordEmail" required />
                        </div>
                        <button type="submit">Request Reset</button>
                    </form>

                    <h3>Reset Password (POST /reset-password)</h3>
                    <form id="resetPasswordForm" onSubmit={handleResetPasswordFormSubmit}>
                        <div className="form-group-item">
                            <label htmlFor="resetToken">Token:</label>
                            <input type="text" id="resetToken" name="resetToken" required />
                        </div>
                        <div className="form-group-item">
                            <label htmlFor="newPassword">New Password:</label>
                            <input type="password" id="newPassword" name="newPassword" required />
                        </div>
                        <button type="submit">Reset Password</button>
                    </form>
                </div>
            )}
        </div>
    );
};

export default Login;
