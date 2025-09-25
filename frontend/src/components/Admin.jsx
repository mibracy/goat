import React, { useState, useEffect } from 'react';
import { callApi } from '../utils/api';

const Admin = ({ activeModel, setApiResponse }) => {
    const [activeAdminSection, setActiveAdminSection] = useState('users');
    const [activeTicketForm, setActiveTicketForm] = useState('create');
    const [ticketDetails, setTicketDetails] = useState(null);
    const [ticketComments, setTicketComments] = useState([]);
    const [commentSearchTerm, setCommentSearchTerm] = useState('');

    const showAdminSection = (sectionName) => {
        setActiveAdminSection(sectionName);
        setApiResponse('');
        setTicketDetails(null); // Clear ticket details when switching sections
        setTicketComments([]); // Clear comments when switching sections
    };

    const showTicketForm = (formType) => {
        setActiveTicketForm(formType);
        setApiResponse('');
        setTicketDetails(null); // Clear ticket details when switching ticket forms
    };

    // Admin User Forms
    const handleListUsers = () => callApi("GET", `/admin/`, null, setApiResponse);
    const handleListAdmins = () => callApi("GET", `/admin/role/Admin`, null, setApiResponse);
    const handleListAgents = () => callApi("GET", `/admin/role/Agent`, null, setApiResponse);
    const handleListCustomers = () => callApi("GET", `/admin/role/Customer`, null, setApiResponse);

    const handleCreateUser = async (e) => {
        e.preventDefault();
        const name = e.target.createName.value;
        const email = e.target.createEmail.value;
        const role = e.target.createRole.value;
        await callApi("POST", `/admin/`, { Name: name, Email: email, Role: role, PasswordHash: "password" }, setApiResponse);
    };

    const handleGetUser = async (e) => {
        e.preventDefault();
        const id = e.target.getUserId.value;
        await callApi("GET", `/admin/${id}`, null, setApiResponse);
    };

    const handleUpdateUser = async (e) => {
        e.preventDefault();
        const id = e.target.updateUserId.value;
        const name = e.target.updateName.value;
        const email = e.target.updateEmail.value;
        const role = e.target.updateRole.value;
        await callApi("PUT", `/admin/update/${id}`, { Name: name, Email: email, Role: role }, setApiResponse);
    };

    const handleDeleteUser = async (e) => {
        e.preventDefault();
        const id = e.target.deleteUserId.value;
        await callApi("DELETE", `/admin/delete/${id}`, null, setApiResponse);
    };

    // Admin Ticket Forms
    const handleListTickets = () => callApi("GET", `/admin/tickets`, null, setApiResponse);

    const handleCreateTicket = async (e) => {
        e.preventDefault();
        const title = e.target.createTicketTitle.value;
        const description = e.target.createTicketDescription.value;
        const status = e.target.createTicketStatus.value;
        const priority = e.target.createTicketPriority.value;
        const requesterId = parseInt(e.target.createTicketRequesterId.value);
        const assigneeIdInput = e.target.createTicketAssigneeId.value;
        const assigneeId = assigneeIdInput ? parseInt(assigneeIdInput) : null;

        const ticketData = {
            Title: title,
            Description: description,
            Status: status,
            Priority: priority,
            RequesterID: requesterId,
        };

        if (assigneeId !== null) {
            ticketData.AssigneeID = assigneeId;
        }

        await callApi("POST", `/admin/tickets`, ticketData, setApiResponse);
    };

    const handleLoadTicket = async (e) => {
        e.preventDefault();
        const ticketId = e.target.updateTicketId.value;
        if (!ticketId) return;

        const data = await callApi("GET", `/admin/tickets/${ticketId}`, null, setApiResponse);
        if (data) {
            setTicketDetails(data);
        } else {
            setTicketDetails(null);
        }
    };

    const handleUpdateTicket = async (e) => {
        e.preventDefault();
        const ticketId = ticketDetails.ID;
        const title = e.target.updateTicketTitle.value;
        const description = e.target.updateTicketDescription.value;
        const status = e.target.updateTicketStatus.value;
        const priority = e.target.updateTicketPriority.value;
        const requesterId = parseInt(e.target.updateTicketRequesterId.value);
        const assigneeIdInput = e.target.updateTicketAssigneeId.value;
        const assigneeId = assigneeIdInput ? parseInt(assigneeIdInput) : null;

        const ticketData = {
            ID: ticketId,
            Title: title,
            Description: description,
            Status: status,
            Priority: priority,
            RequesterID: requesterId,
        };

        if (assigneeId !== null) {
            ticketData.AssigneeID = assigneeId;
        }

        await callApi("PUT", `/admin/tickets/${ticketId}`, ticketData, setApiResponse);
    };

    // Admin Comment Forms
    const handleListComments = () => callApi("GET", `/admin/comments`, null, setApiResponse);

    const handleLoadTicketDetails = async (e) => {
        e.preventDefault();
        const ticketId = e.target.viewTicketId.value;
        if (!ticketId) return;

        const data = await callApi("GET", `/admin/tickets/${ticketId}`, null, setApiResponse);
        if (data) {
            setTicketDetails(data);
            setTicketComments(data.Comments || []);
        } else {
            setTicketDetails(null);
            setTicketComments([]);
        }
    };

    const handleAddComment = async (e) => {
        e.preventDefault();
        const ticketId = parseInt(e.target.addCommentTicketId.value);
        const authorId = parseInt(e.target.addCommentAuthorId.value);
        const body = e.target.addCommentBody.value;
        const isInternal = e.target.addCommentIsInternal.checked;

        const commentData = {
            TicketID: ticketId,
            AuthorID: authorId,
            Body: body,
            IsInternal: isInternal,
        };

        await callApi("POST", `/admin/comments`, commentData, setApiResponse);
        // Re-fetch ticket details to update comments section
        if (ticketDetails) {
            handleLoadTicketDetails({ preventDefault: () => {}, target: { viewTicketId: { value: ticketDetails.ID } } });
        }
    };

    const filteredComments = ticketComments.filter(comment =>
        comment.Body.toLowerCase().includes(commentSearchTerm.toLowerCase())
    );

    useEffect(() => {
        // Enforce numeric input for elements with class 'js-number-input' and handle arrow key increments/decrements
        document.querySelectorAll(".js-number-input").forEach((inputElement) => {
            inputElement.addEventListener("input", function (event) {
                this.value = this.value.replace(/[^0-9]/g, "");
            });

            inputElement.addEventListener("keydown", function (event) {
                let currentValue = parseInt(this.value) || 0;
                if (event.key === "ArrowUp") {
                    event.preventDefault();
                    this.value = currentValue + 1;
                } else if (event.key === "ArrowDown") {
                    event.preventDefault();
                    this.value = Math.max(0, currentValue - 1);
                }
            });
        });
    }, [activeAdminSection]);

    if (activeModel !== 'admin') {
        return null;
    }

    return (
        <div id="admin" className="model-section active">
            <h2>Admin</h2>
            <div className="menu">
                <button onClick={() => showAdminSection('users')} className={activeAdminSection === 'users' ? 'active' : ''}>
                    Users
                </button>
                <button onClick={() => showAdminSection('tickets')} className={activeAdminSection === 'tickets' ? 'active' : ''}>Tickets</button>
                <button onClick={() => showAdminSection('comments')} className={activeAdminSection === 'comments' ? 'active' : ''}>Comments</button>
            </div>

            {activeAdminSection === 'users' && (
                <div id="admin-users" className="admin-sub-section active">
                    <h3>List Users (GET /admin/) (GET /admin/role/{role})</h3>
                    <div className="button-group">
                        <form onSubmit={handleListUsers}>
                            <button type="submit">All</button>
                        </form>
                        <form onSubmit={handleListAdmins}>
                            <button type="submit">Admins</button>
                        </form>
                        <form onSubmit={handleListAgents}>
                            <button type="submit">Agents</button>
                        </form>
                        <form onSubmit={handleListCustomers}>
                            <button type="submit">Customers</button>
                        </form>
                    </div>

                    <div id="admin-create-user" className="disabled-admin-sub-section">
                        <h3>Create User (POST /admin/)</h3>
                        <form onSubmit={handleCreateUser}>
                            <div className="form-group-item">
                                <label htmlFor="createName">Name:</label>
                                <input type="text" id="createName" name="createName" required />
                            </div>
                            <div className="form-group-item">
                                <label htmlFor="createEmail">Email:</label>
                                <input type="text" id="createEmail" name="createEmail" required />
                            </div>
                            <div className="form-group-item">
                                <label htmlFor="createRole">Role:</label>
                                <select id="createRole" name="createRole">
                                    <option value="Customer">Customer</option>
                                    <option value="Agent">Agent</option>
                                    <option value="Admin">Admin</option>
                                </select>
                            </div>
                            <button type="submit">Create User</button>
                        </form>
                    </div>

                    <div id="admin-get-user" className="disabled-admin-sub-section">
                        <h3>Get User by ID (GET /admin/{id})</h3>
                        <form onSubmit={handleGetUser}>
                            <label htmlFor="getUserId">User ID:</label>
                            <input type="text" id="getUserId" name="getUserId" required className="js-number-input" />
                            <button type="submit">Get User</button>
                        </form>
                    </div>

                    <div id="admin-update-user" className="disabled-admin-sub-section">
                        <h3>Update User (PUT /admin/update/{id})</h3>
                        <form onSubmit={handleUpdateUser}>
                            <div className="form-group-item">
                                <label htmlFor="updateUserId">User ID:</label>
                                <input type="text" id="updateUserId" name="updateUserId" required className="js-number-input" />
                            </div>
                            <div className="form-group-item">
                                <label htmlFor="updateName">New Name:</label>
                                <input type="text" id="updateName" name="updateName" required />
                            </div>
                            <div className="form-group-item">
                                <label htmlFor="updateEmail">New Email:</label>
                                <input type="text" id="updateEmail" name="updateEmail" required />
                            </div>
                            <div className="form-group-item">
                                <label htmlFor="updateRole">New Role:</label>
                                <select id="updateRole" name="updateRole">
                                    <option value="Customer">Customer</option>
                                    <option value="Agent">Agent</option>
                                    <option value="Admin">Admin</option>
                                </select>
                            </div>
                            <button type="submit">Update User</button>
                        </form>
                    </div>

                    <div id="admin-delete-user" className="disabled-admin-sub-section">
                        <h3>Delete User (DELETE /admin/delete/{id})</h3>
                        <form onSubmit={handleDeleteUser}>
                            <label htmlFor="deleteUserId">User ID:</label>
                            <input type="text" id="deleteUserId" name="deleteUserId" required className="js-number-input" />
                            <button type="submit">Delete User</button>
                        </form>
                    </div>
                </div>
            )}

            {activeAdminSection === 'tickets' && (
                <div id="admin-tickets" className="admin-sub-section">
                    <h2>Tickets</h2>
                    <div>
                        <h3>List Tickets (GET /admin/tickets)</h3>
                        <form onSubmit={handleListTickets}>
                            <button type="submit">List All Tickets</button>
                        </form>
                    </div>

                    <div>
                        <h3>Create/Update Ticket</h3>
                        <div className="button-group">
                            <button onClick={() => showTicketForm('create')} className={activeTicketForm === 'create' ? 'active' : ''}>
                                Create Ticket
                            </button>
                            <button onClick={() => showTicketForm('update')} className={activeTicketForm === 'update' ? 'active' : ''}>
                                Update Ticket
                            </button>
                        </div>

                        {activeTicketForm === 'create' && (
                            <div id="createTicketSection">
                                <form onSubmit={handleCreateTicket}>
                                    <div className="form-group-item">
                                        <label htmlFor="createTicketTitle">Title:</label>
                                        <input type="text" id="createTicketTitle" name="createTicketTitle" required />
                                    </div>

                                    <div className="form-group-item full-width-item">
                                        <label htmlFor="createTicketDescription">Description:</label>
                                        <textarea id="createTicketDescription" name="createTicketDescription"></textarea>
                                    </div>

                                    <div className="form-group-grid">
                                        <div className="form-group-item">
                                            <label htmlFor="createTicketStatus">Status:</label>
                                            <select id="createTicketStatus" name="createTicketStatus">
                                                <option value="Open">Open</option>
                                            </select>
                                        </div>
                                        <div className="form-group-item">
                                            <label htmlFor="createTicketPriority">Priority:</label>
                                            <select id="createTicketPriority" name="createTicketPriority">
                                                <option value="Low">Low</option>
                                                <option value="Medium">Medium</option>
                                                <option value="High">High</option>
                                                <option value="Urgent">Urgent</option>
                                            </select>
                                        </div>
                                    </div>

                                    <div className="form-group-grid">
                                        <div className="form-group-item">
                                            <label htmlFor="createTicketRequesterId">Requester ID:</label>
                                            <input type="text" id="createTicketRequesterId" name="createTicketRequesterId" required className="js-number-input" />
                                        </div>
                                        <div className="form-group-item">
                                            <label htmlFor="createTicketAssigneeId">Assignee ID (Optional):</label>
                                            <input type="text" id="createTicketAssigneeId" name="createTicketAssigneeId" className="js-number-input" />
                                        </div>
                                    </div>

                                    <button type="submit">Create Ticket</button>
                                </form>
                            </div>
                        )}

                        {activeTicketForm === 'update' && (
                            <div id="updateTicketSection">
                                <form onSubmit={handleLoadTicket}>
                                    <label htmlFor="updateTicketId">Ticket ID:</label>
                                    <input type="text" id="updateTicketId" name="updateTicketId" required className="js-number-input" />
                                    <button type="submit">Load Ticket</button>
                                </form>
                                {ticketDetails && (
                                    <form id="updateTicketForm" style={{ paddingTop: '5px', display: 'grid' }} onSubmit={handleUpdateTicket}>
                                        <input type="hidden" id="hiddenUpdateTicketId" name="hiddenUpdateTicketId" value={ticketDetails.ID} />
                                        <div className="form-group-item">
                                            <label htmlFor="updateTicketTitle">Title:</label>
                                            <input type="text" id="updateTicketTitle" name="updateTicketTitle" required value={ticketDetails.Title} onChange={(e) => setTicketDetails({ ...ticketDetails, Title: e.target.value })} />
                                        </div>

                                        <div className="form-group-item full-width-item">
                                            <label htmlFor="updateTicketDescription">Description:</label>
                                            <textarea id="updateTicketDescription" name="updateTicketDescription" value={ticketDetails.Description} onChange={(e) => setTicketDetails({ ...ticketDetails, Description: e.target.value })}></textarea>
                                        </div>

                                        <div className="form-group-grid">
                                            <div className="form-group-item">
                                                <label htmlFor="updateTicketStatus">Status:</label>
                                                <select id="updateTicketStatus" name="updateTicketStatus" value={ticketDetails.Status} onChange={(e) => setTicketDetails({ ...ticketDetails, Status: e.target.value })}>
                                                    <option value="Open">Open</option>
                                                    <option value="In Progress">In Progress</option>
                                                    <option value="Resolved">Resolved</option>
                                                    <option value="Closed">Closed</option>
                                                </select>
                                            </div>
                                            <div className="form-group-item">
                                                <label htmlFor="updateTicketPriority">Priority:</label>
                                                <select id="updateTicketPriority" name="updateTicketPriority" value={ticketDetails.Priority} onChange={(e) => setTicketDetails({ ...ticketDetails, Priority: e.target.value })}>
                                                    <option value="Low">Low</option>
                                                    <option value="Medium">Medium</option>
                                                    <option value="High">High</option>
                                                    <option value="Urgent">Urgent</option>
                                                </select>
                                            </div>
                                        </div>

                                        <div className="form-group-grid">
                                            <div className="form-group-item">
                                                <label htmlFor="updateTicketRequesterId">Requester ID:</label>
                                                <input type="text" id="updateTicketRequesterId" name="updateTicketRequesterId" required value={ticketDetails.RequesterID} onChange={(e) => setTicketDetails({ ...ticketDetails, RequesterID: parseInt(e.target.value) })} className="js-number-input" />
                                            </div>
                                            <div className="form-group-item">
                                                <label htmlFor="updateTicketAssigneeId">Assignee ID (Optional):</label>
                                                <input type="text" id="updateTicketAssigneeId" name="updateTicketAssigneeId" value={ticketDetails.AssigneeID.Int64 || ''} onChange={(e) => setTicketDetails({ ...ticketDetails, AssigneeID: { Int64: parseInt(e.target.value), Valid: true } })} className="js-number-input" />
                                            </div>
                                        </div>

                                        <button type="submit">Update Ticket</button>
                                    </form>
                                )}
                            </div>
                        )}
                    </div>
                </div>
            )}

            {activeAdminSection === 'comments' && (
                <div id="admin-comments" className="admin-sub-section">
                    <h2>Comments</h2>
                    <div className="button-group">
                        <button onClick={handleListComments} className="active">
                            List All Comments
                        </button>
                        <button onClick={() => showAdminSection('comments')}> {/* This button now just re-renders the comments section */}
                            Lookup Ticket
                        </button>
                    </div>

                    <div id="viewTicketDetailsSection" style={{ display: 'block' }}>
                        <h3>Ticket Details</h3>
                        <form onSubmit={handleLoadTicketDetails}>
                            <label htmlFor="viewTicketId">Ticket ID:</label>
                            <input type="text" id="viewTicketId" name="viewTicketId" required className="js-number-input" />
                            <button type="submit">Load Ticket Details</button>
                        </form>
                        {ticketDetails && (
                            <form id="ticketDetailsForm" style={{ paddingTop: '5px', display: 'grid' }}>
                                <div className="form-group-item">
                                    <label>ID:</label>
                                    <input type="text" readOnly value={ticketDetails.ID} />
                                </div>
                                <div className="form-group-item">
                                    <label>Title:</label>
                                    <input type="text" readOnly value={ticketDetails.Title} />
                                </div>
                                <div className="form-group-item full-width-item">
                                    <label>Description:</label>
                                    <textarea readOnly value={ticketDetails.Description}></textarea>
                                </div>

                                <div className="form-group-grid">
                                    <div className="form-group-item">
                                        <label>Status:</label>
                                        <input type="text" readOnly value={ticketDetails.Status} />
                                    </div>
                                    <div className="form-group-item">
                                        <label>Priority:</label>
                                        <input type="text" readOnly value={ticketDetails.Priority} />
                                    </div>
                                </div>

                                <div className="form-group-grid">
                                    <div className="form-group-item">
                                        <label>Requester ID:</label>
                                        <input type="text" readOnly value={ticketDetails.RequesterID} />
                                    </div>
                                    <div className="form-group-item">
                                        <label>Assignee ID:</label>
                                        <input type="text" readOnly value={ticketDetails.AssigneeID.Int64 || ''} />
                                    </div>
                                </div>

                                <div className="form-group-grid">
                                    <div className="form-group-item">
                                        <label>Created At:</label>
                                        <input type="text" readOnly value={new Date(ticketDetails.CreatedAt).toLocaleString()} />
                                    </div>
                                    <div className="form-group-item">
                                        <label>Updated At:</label>
                                        <input type="text" readOnly value={new Date(ticketDetails.UpdatedAt).toLocaleString()} />
                                    </div>
                                    <div className="form-group-item">
                                        <label>Closed At:</label>
                                        <input type="text" readOnly value={ticketDetails.ClosedAt.Valid ? new Date(ticketDetails.ClosedAt.Time).toLocaleString() : ''} />
                                    </div>
                                </div>
                            </form>
                        )}

                        <div id="addCommentSection" style={{ display: ticketDetails ? 'block' : 'none' }}>
                            <h3>Add New Comment</h3>
                            <form onSubmit={handleAddComment}>
                                <div className="form-group-item">
                                    <label htmlFor="addCommentTicketId">Ticket ID:</label>
                                    <input type="text" id="addCommentTicketId" name="addCommentTicketId" required readOnly value={ticketDetails?.ID || ''} />
                                </div>
                                <div className="form-group-item">
                                    <label htmlFor="addCommentAuthorId">Author ID:</label>
                                    <input type="text" id="addCommentAuthorId" name="addCommentAuthorId" required className="js-number-input" />
                                </div>
                                <div className="form-group-item full-width-item">
                                    <label htmlFor="addCommentBody">Comment:</label>
                                    <textarea id="addCommentBody" name="addCommentBody" required></textarea>
                                </div>
                                <div className="form-group-item">
                                    <label htmlFor="addCommentIsInternal">Internal Only:</label>
                                    <input type="checkbox" id="addCommentIsInternal" name="addCommentIsInternal" style={{ maxWidth: '25%', minHeight: '42px' }} />
                                </div>
                                <button type="submit">Add Comment</button>
                            </form>
                        </div>

                        <div id="ticketCommentsSection" style={{ display: ticketDetails ? 'block' : 'none' }}>
                            <h3>Comments for this Ticket</h3>
                            <input type="text" id="commentSearchInput" placeholder="Search comments..." style={{ width: '90%', marginBottom: '10px', marginLeft: '3%' }} onChange={(e) => setCommentSearchTerm(e.target.value)} />
                            <div style={{ maxHeight: '600px', overflowY: 'auto', border: '1px solid #0f0' }}>
                                <table id="ticketCommentsTable" style={{ width: '100%', marginTop: '10px' }}>
                                    <thead>
                                        <tr>
                                            <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'center', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>ID</th>
                                            <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'center', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>Author ID</th>
                                            <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'left', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>Body</th>
                                            <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'center', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>Internal</th>
                                            <th style={{ border: '1px solid #0f0', padding: '8px', textAlign: 'left', position: 'sticky', top: 0, backgroundColor: '#000', zIndex: 1 }}>Created At</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {filteredComments.map(comment => (
                                            <tr key={comment.ID}>
                                                <td style={{ border: '1px solid #0f0', textAlign: 'center' }}>{comment.ID}</td>
                                                <td style={{ border: '1px solid #0f0', textAlign: 'center' }}>{comment.AuthorID}</td>
                                                <td style={{ border: '1px solid #0f0' }}>{comment.Body}</td>
                                                <td style={{ border: '1px solid #0f0', textAlign: 'center' }}>{comment.IsInternal ? "Yes" : "No"}</td>
                                                <td style={{ border: '1px solid #0f0', textAlign: 'left' }}>{new Date(comment.CreatedAt).toLocaleString()}</td>
                                            </tr>
                                        ))}
                                        {filteredComments.length === 0 && (
                                            <tr>
                                                <td colSpan="5" style={{ border: '1px solid #0f0' }}>No comments found for this ticket.</td>
                                            </tr>
                                        )}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Admin;