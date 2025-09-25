import React, { useState } from 'react';

const Layout = ({ children }) => {
    const [activeModel, setActiveModel] = useState('login');
    const [apiResponse, setApiResponse] = useState('');

    const showModel = (modelName) => {
        setActiveModel(modelName);
        setApiResponse(''); // Clear API response when switching models
    };

    return (
        <>
            <h1>GOAT API Viewer</h1>

            <div className="menu">
                <button onClick={() => showModel('login')} className={activeModel === 'login' ? 'active' : ''}>Login</button>
                <button onClick={() => showModel('admin')} className={activeModel === 'admin' ? 'active' : ''}>Admin</button>
                <button onClick={() => showModel('agent')} className={activeModel === 'agent' ? 'active' : ''}>Agent</button>
                <button onClick={() => showModel('customer')} className={activeModel === 'customer' ? 'active' : ''}>Customer</button>
            </div>

            {React.Children.map(children, child => {
                if (React.isValidElement(child)) {
                    return React.cloneElement(child, { activeModel, setApiResponse });
                }
                return child;
            })}

            <h2>API Response:</h2>
            <div id="response-container">
                <pre id="apiResponse">{apiResponse}</pre>
            </div>
        </>
    );
};

export default Layout;
