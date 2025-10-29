import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function App() {
  const [users, setUsers] = useState([]);
  const [tasks, setTasks] = useState([]);
  const [activeTab, setActiveTab] = useState('dashboard');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [currentUser, setCurrentUser] = useState(null);
  const [showLogin, setShowLogin] = useState(true);
  const [showRegister, setShowRegister] = useState(false);

  // Состояния для форм
  const [newUser, setNewUser] = useState({
    username: '',
    email: '',
    first_name: '',
    last_name: '',
    role: 'user'
  });
  const [newTask, setNewTask] = useState({
    title: '',
    description: '',
    status: 'pending',
    priority: 'medium',
    assigned_to: '',
    project_id: 1
  });
  const [editingUser, setEditingUser] = useState(null);
  const [editingTask, setEditingTask] = useState(null);
  const [showUserForm, setShowUserForm] = useState(false);
  const [showTaskForm, setShowTaskForm] = useState(false);

  // Состояния для логина и регистрации
  const [loginData, setLoginData] = useState({
    username: '',
    password: ''
  });
  const [registerData, setRegisterData] = useState({
    username: '',
    email: '',
    password: '',
    first_name: '',
    last_name: ''
  });

  useEffect(() => {
    if (currentUser) {
      fetchUsers();
      fetchTasks();
    }
  }, [currentUser, activeTab]);

  const showMessage = (message, type = 'success') => {
    if (type === 'success') {
      setSuccess(message);
      setTimeout(() => setSuccess(''), 3000);
    } else {
      setError(message);
      setTimeout(() => setError(''), 5000);
    }
  };

  const handleLogin = (e) => {
    e.preventDefault();
    // Простая имитация авторизации
    const user = {
      id: 1,
      username: loginData.username,
      email: `${loginData.username}@company.com`,
      first_name: 'Пользователь',
      last_name: 'Системы',
      role: loginData.username === 'admin' ? 'admin' : 
            loginData.username === 'manager' ? 'manager' : 'user'
    };
    setCurrentUser(user);
    setShowLogin(false);
    showMessage(`Добро пожаловать, ${user.username}!`);
  };

  const handleRegister = async (e) => {
    e.preventDefault();
    try {
      await axios.post(`${API_URL}/users`, {
        username: registerData.username,
        email: registerData.email,
        first_name: registerData.first_name,
        last_name: registerData.last_name,
        role: 'user'
      });
      showMessage('✅ Регистрация успешна! Теперь вы можете войти.');
      setShowRegister(false);
      setRegisterData({ username: '', email: '', password: '', first_name: '', last_name: '' });
    } catch (error) {
      console.error('Error registering user:', error);
      showMessage('❌ Ошибка регистрации', 'error');
    }
  };

  const handleLogout = () => {
    setCurrentUser(null);
    setShowLogin(true);
    setUsers([]);
    setTasks([]);
    showMessage('Вы вышли из системы');
  };

  const fetchUsers = async () => {
    if (!currentUser) return;
    
    setLoading(true);
    setError('');
    try {
      const response = await axios.get(`${API_URL}/users`);
      const usersData = response.data.users || response.data || [];
      setUsers(Array.isArray(usersData) ? usersData : []);
    } catch (error) {
      console.error('Error fetching users:', error);
      showMessage('Ошибка при загрузке пользователей', 'error');
    } finally {
      setLoading(false);
    }
  };

  const fetchTasks = async () => {
    if (!currentUser) return;
    
    setLoading(true);
    setError('');
    try {
      const response = await axios.get(`${API_URL}/tasks`);
      const tasksData = response.data.tasks || response.data || [];
      setTasks(Array.isArray(tasksData) ? tasksData : []);
    } catch (error) {
      console.error('Error fetching tasks:', error);
      showMessage('Ошибка при загрузке задач', 'error');
    } finally {
      setLoading(false);
    }
  };

  const createUser = async (e) => {
    e.preventDefault();
    try {
      await axios.post(`${API_URL}/users`, newUser);
      setNewUser({ username: '', email: '', first_name: '', last_name: '', role: 'user' });
      setShowUserForm(false);
      fetchUsers();
      showMessage('👤 Пользователь успешно создан!');
    } catch (error) {
      console.error('Error creating user:', error);
      showMessage('❌ Ошибка при создании пользователя', 'error');
    }
  };

  const createTask = async (e) => {
    e.preventDefault();
    try {
      const taskData = {
        ...newTask,
        created_by: currentUser.id,
        assigned_to: newTask.assigned_to || currentUser.id
      };
      await axios.post(`${API_URL}/tasks`, taskData);
      setNewTask({ title: '', description: '', status: 'pending', priority: 'medium', assigned_to: '', project_id: 1 });
      setShowTaskForm(false);
      fetchTasks();
      showMessage('✅ Задача успешно создана!');
    } catch (error) {
      console.error('Error creating task:', error);
      showMessage('❌ Ошибка при создании задачи', 'error');
    }
  };

  const updateUser = async (e) => {
    e.preventDefault();
    try {
      await axios.put(`${API_URL}/users/${editingUser.id}`, editingUser);
      setEditingUser(null);
      fetchUsers();
      showMessage('👤 Пользователь успешно обновлен!');
    } catch (error) {
      console.error('Error updating user:', error);
      showMessage('❌ Ошибка при обновлении пользователя', 'error');
    }
  };

  const updateTask = async (e) => {
    e.preventDefault();
    try {
      await axios.put(`${API_URL}/tasks/${editingTask.id}`, editingTask);
      setEditingTask(null);
      fetchTasks();
      showMessage('✅ Задача успешно обновлена!');
    } catch (error) {
      console.error('Error updating task:', error);
      showMessage('❌ Ошибка при обновлении задачи', 'error');
    }
  };

  const deleteUser = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить пользователя?')) {
      try {
        await axios.delete(`${API_URL}/users/${id}`);
        fetchUsers();
        showMessage('👤 Пользователь удален!');
      } catch (error) {
        console.error('Error deleting user:', error);
        showMessage('❌ Ошибка при удалении пользователя', 'error');
      }
    }
  };

  const deleteTask = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить задачу?')) {
      try {
        await axios.delete(`${API_URL}/tasks/${id}`);
        fetchTasks();
        showMessage('✅ Задача удалена!');
      } catch (error) {
        console.error('Error deleting task:', error);
        showMessage('❌ Ошибка при удалении задачи', 'error');
      }
    }
  };

  const openEditUser = (user) => {
    setEditingUser({...user});
  };

  const openEditTask = (task) => {
    setEditingTask({...task});
  };

  const getUserName = (userId) => {
    const user = users.find(u => u.id === userId);
    return user ? user.username : 'Не назначена';
  };

  const canAssignToUsers = currentUser?.role === 'admin' || currentUser?.role === 'manager';

  if (showLogin) {
    return (
      <div className="login-container">
        <div className="login-form">
          <h1>🚀 Микросервисная система</h1>
          <p>Войдите в систему для продолжения</p>
          
          {showRegister ? (
            <form onSubmit={handleRegister}>
              <h3>📝 Регистрация</h3>
              <div className="form-group">
                <label>Имя пользователя *</label>
                <input
                  type="text"
                  value={registerData.username}
                  onChange={(e) => setRegisterData({...registerData, username: e.target.value})}
                  placeholder="Введите имя пользователя"
                  required
                />
              </div>
              <div className="form-group">
                <label>Email *</label>
                <input
                  type="email"
                  value={registerData.email}
                  onChange={(e) => setRegisterData({...registerData, email: e.target.value})}
                  placeholder="Введите email"
                  required
                />
              </div>
              <div className="form-group">
                <label>Пароль *</label>
                <input
                  type="password"
                  value={registerData.password}
                  onChange={(e) => setRegisterData({...registerData, password: e.target.value})}
                  placeholder="Введите пароль"
                  required
                />
              </div>
              <div className="form-group">
                <label>Имя</label>
                <input
                  type="text"
                  value={registerData.first_name}
                  onChange={(e) => setRegisterData({...registerData, first_name: e.target.value})}
                  placeholder="Введите имя"
                />
              </div>
              <div className="form-group">
                <label>Фамилия</label>
                <input
                  type="text"
                  value={registerData.last_name}
                  onChange={(e) => setRegisterData({...registerData, last_name: e.target.value})}
                  placeholder="Введите фамилию"
                />
              </div>
              <button type="submit" className="register-btn">📝 Зарегистрироваться</button>
              <button 
                type="button" 
                className="back-btn"
                onClick={() => setShowRegister(false)}
              >
                ← Назад к входу
              </button>
            </form>
          ) : (
            <form onSubmit={handleLogin}>
              <div className="form-group">
                <label>Имя пользователя:</label>
                <input
                  type="text"
                  value={loginData.username}
                  onChange={(e) => setLoginData({...loginData, username: e.target.value})}
                  placeholder="Введите имя пользователя"
                  required
                />
              </div>
              <div className="form-group">
                <label>Пароль:</label>
                <input
                  type="password"
                  value={loginData.password}
                  onChange={(e) => setLoginData({...loginData, password: e.target.value})}
                  placeholder="Введите пароль"
                  required
                />
              </div>
              <button type="submit" className="login-btn">🔐 Войти в систему</button>
              <button 
                type="button" 
                className="register-btn"
                onClick={() => setShowRegister(true)}
              >
                📝 Создать аккаунт
              </button>
            </form>
          )}

          <div className="login-hint">
            <p><strong>Тестовые пользователи:</strong></p>
            <p>👤 admin / password123 (Администратор)</p>
            <p>👤 manager / password123 (Менеджер)</p>
            <p>👤 user1 / password123 (Пользователь)</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="App">
      <header className="App-header">
        <div className="header-content">
          <div className="header-info">
            <h1>🚀 Микросервисная система управления</h1>
            <div className="user-info">
              <span>👋 Привет, {currentUser?.username}!</span>
              <span className={`user-role role-${currentUser?.role}`}>
                {currentUser?.role === 'admin' ? '🔧 Администратор' : 
                 currentUser?.role === 'manager' ? '👔 Менеджер' : '👤 Пользователь'}
              </span>
            </div>
          </div>
          
          <nav className="nav-tabs">
            <button 
              onClick={() => setActiveTab('dashboard')}
              className={`tab-button ${activeTab === 'dashboard' ? 'active' : ''}`}
            >
              📊 Дашборд
            </button>
            <button 
              onClick={() => setActiveTab('users')}
              className={`tab-button ${activeTab === 'users' ? 'active' : ''}`}
            >
              👥 Пользователи ({users.length})
            </button>
            <button 
              onClick={() => setActiveTab('tasks')}
              className={`tab-button ${activeTab === 'tasks' ? 'active' : ''}`}
            >
              ✅ Задачи ({tasks.length})
            </button>
            <button 
              onClick={handleLogout}
              className="logout-btn"
            >
              🚪 Выйти
            </button>
          </nav>
        </div>
      </header>

      <main className="main-content">
        {error && (
          <div className="message error">
            ⚠️ {error}
          </div>
        )}

        {success && (
          <div className="message success">
            ✅ {success}
          </div>
        )}

        {loading && (
          <div className="loading">
            ⏳ Загрузка...
          </div>
        )}

        {activeTab === 'dashboard' && (
          <div className="dashboard">
            <h2>📊 Дашборд системы</h2>
            <div className="stats-grid">
              <div className="stat-card">
                <div className="stat-icon">👥</div>
                <div className="stat-info">
                  <h3>Пользователи</h3>
                  <div className="stat-number">{users.length}</div>
                </div>
                <button className="refresh-btn" onClick={fetchUsers} title="Обновить">
                  🆙
                </button>
              </div>
              
              <div className="stat-card">
                <div className="stat-icon">✅</div>
                <div className="stat-info">
                  <h3>Задачи</h3>
                  <div className="stat-number">{tasks.length}</div>
                </div>
                <button className="refresh-btn" onClick={fetchTasks} title="Обновить">
                  🆙
                </button>
              </div>

              <div className="stat-card">
                <div className="stat-icon">⚡</div>
                <div className="stat-info">
                  <h3>В работе</h3>
                  <div className="stat-number">
                    {tasks.filter(task => task.status === 'in_progress').length}
                  </div>
                </div>
              </div>

              <div className="stat-card">
                <div className="stat-icon">🎯</div>
                <div className="stat-info">
                  <h3>Завершено</h3>
                  <div className="stat-number">
                    {tasks.filter(task => task.status === 'completed').length}
                  </div>
                </div>
              </div>
            </div>

            <div className="recent-activity">
              <h3>📈 Быстрый обзор</h3>
              <div className="activity-grid">
                <div className="activity-card">
                  <h4>👥 Распределение по ролям</h4>
                  <div className="role-stats">
                    <div>🔧 Администраторы: {users.filter(u => u.role === 'admin').length}</div>
                    <div>👔 Менеджеры: {users.filter(u => u.role === 'manager').length}</div>
                    <div>👤 Пользователи: {users.filter(u => u.role === 'user').length}</div>
                  </div>
                </div>
                <div className="activity-card">
                  <h4>✅ Статусы задач</h4>
                  <div className="task-stats">
                    <div>⏳ Ожидает: {tasks.filter(t => t.status === 'pending').length}</div>
                    <div>⚡ В работе: {tasks.filter(t => t.status === 'in_progress').length}</div>
                    <div>✅ Завершено: {tasks.filter(t => t.status === 'completed').length}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'users' && (
          <div className="users-section">
            <div className="section-header">
              <h2>👥 Управление пользователями</h2>
              <div className="action-buttons">
                <button className="refresh-btn large" onClick={fetchUsers} title="Обновить список">
                  🆙 Обновить
                </button>
                {currentUser?.role === 'admin' && (
                  <button 
                    className="add-btn"
                    onClick={() => setShowUserForm(!showUserForm)}
                  >
                    ➕ Добавить пользователя
                  </button>
                )}
              </div>
            </div>

            {showUserForm && (
              <form className="form" onSubmit={createUser}>
                <h3>➕ Новый пользователь</h3>
                <div className="form-grid">
                  <div className="form-group">
                    <label>Имя пользователя *</label>
                    <input
                      type="text"
                      value={newUser.username}
                      onChange={(e) => setNewUser({...newUser, username: e.target.value})}
                      placeholder="Введите имя пользователя"
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>Email *</label>
                    <input
                      type="email"
                      value={newUser.email}
                      onChange={(e) => setNewUser({...newUser, email: e.target.value})}
                      placeholder="Введите email"
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>Имя</label>
                    <input
                      type="text"
                      value={newUser.first_name}
                      onChange={(e) => setNewUser({...newUser, first_name: e.target.value})}
                      placeholder="Введите имя"
                    />
                  </div>
                  <div className="form-group">
                    <label>Фамилия</label>
                    <input
                      type="text"
                      value={newUser.last_name}
                      onChange={(e) => setNewUser({...newUser, last_name: e.target.value})}
                      placeholder="Введите фамилию"
                    />
                  </div>
                  <div className="form-group">
                    <label>Роль</label>
                    <select
                      value={newUser.role}
                      onChange={(e) => setNewUser({...newUser, role: e.target.value})}
                    >
                      <option value="user">👤 Пользователь</option>
                      <option value="manager">👔 Менеджер</option>
                      <option value="admin">🔧 Администратор</option>
                    </select>
                  </div>
                </div>
                <div className="form-actions">
                  <button type="submit" className="save-btn">💾 Сохранить</button>
                  <button type="button" className="cancel-btn" onClick={() => setShowUserForm(false)}>
                    ❌ Отмена
                  </button>
                </div>
              </form>
            )}

            {editingUser && (
              <form className="form" onSubmit={updateUser}>
                <h3>✏️ Редактирование пользователя</h3>
                <div className="form-grid">
                  <div className="form-group">
                    <label>Имя пользователя *</label>
                    <input
                      type="text"
                      value={editingUser.username}
                      onChange={(e) => setEditingUser({...editingUser, username: e.target.value})}
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>Email *</label>
                    <input
                      type="email"
                      value={editingUser.email}
                      onChange={(e) => setEditingUser({...editingUser, email: e.target.value})}
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>Имя</label>
                    <input
                      type="text"
                      value={editingUser.first_name || ''}
                      onChange={(e) => setEditingUser({...editingUser, first_name: e.target.value})}
                    />
                  </div>
                  <div className="form-group">
                    <label>Фамилия</label>
                    <input
                      type="text"
                      value={editingUser.last_name || ''}
                      onChange={(e) => setEditingUser({...editingUser, last_name: e.target.value})}
                    />
                  </div>
                  <div className="form-group">
                    <label>Роль</label>
                    <select
                      value={editingUser.role}
                      onChange={(e) => setEditingUser({...editingUser, role: e.target.value})}
                    >
                      <option value="user">👤 Пользователь</option>
                      <option value="manager">👔 Менеджер</option>
                      <option value="admin">🔧 Администратор</option>
                    </select>
                  </div>
                </div>
                <div className="form-actions">
                  <button type="submit" className="save-btn">💾 Сохранить</button>
                  <button type="button" className="cancel-btn" onClick={() => setEditingUser(null)}>
                    ❌ Отмена
                  </button>
                </div>
              </form>
            )}

            <div className="users-list">
              {users.length === 0 ? (
                <div className="empty-state">
                  👥 Нет пользователей
                </div>
              ) : (
                <div className="cards-grid">
                  {users.map(user => (
                    <div key={user.id} className="user-card">
                      <div className="card-header">
                        <div className="user-avatar">
                          {user.role === 'admin' ? '🔧' : user.role === 'manager' ? '👔' : '👤'}
                        </div>
                        <div className="user-info">
                          <div className="user-name">{user.username}</div>
                          <div className="user-email">{user.email}</div>
                        </div>
                        <div className="card-actions">
                          {currentUser?.role === 'admin' && (
                            <>
                              <button 
                                className="edit-btn"
                                onClick={() => openEditUser(user)}
                                title="Редактировать"
                              >
                                ✏️
                              </button>
                              <button 
                                className="delete-btn"
                                onClick={() => deleteUser(user.id)}
                                title="Удалить"
                              >
                                🗑️
                              </button>
                            </>
                          )}
                        </div>
                      </div>
                      <div className="card-content">
                        <div className="user-details">
                          {user.first_name && (
                            <div>👤 {user.first_name} {user.last_name}</div>
                          )}
                          <div className={`role-badge role-${user.role}`}>
                            {user.role === 'admin' ? 'Администратор' : 
                             user.role === 'manager' ? 'Менеджер' : 'Пользователь'}
                          </div>
                        </div>
                        <div className="card-footer">
                          <small>Создан: {new Date(user.created_at).toLocaleDateString('ru-RU')}</small>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}

        {activeTab === 'tasks' && (
          <div className="tasks-section">
            <div className="section-header">
              <h2>✅ Управление задачами</h2>
              <div className="action-buttons">
                <button className="refresh-btn large" onClick={fetchTasks} title="Обновить список">
                  🆙 Обновить
                </button>
                <button 
                  className="add-btn"
                  onClick={() => setShowTaskForm(!showTaskForm)}
                >
                  ➕ Добавить задачу
                </button>
              </div>
            </div>

            {showTaskForm && (
              <form className="form" onSubmit={createTask}>
                <h3>➕ Новая задача</h3>
                <div className="form-grid">
                  <div className="form-group full-width">
                    <label>Название задачи *</label>
                    <input
                      type="text"
                      value={newTask.title}
                      onChange={(e) => setNewTask({...newTask, title: e.target.value})}
                      placeholder="Введите название задачи"
                      required
                    />
                  </div>
                  <div className="form-group full-width">
                    <label>Описание</label>
                    <textarea
                      value={newTask.description}
                      onChange={(e) => setNewTask({...newTask, description: e.target.value})}
                      placeholder="Опишите задачу подробнее"
                      rows="3"
                    />
                  </div>
                  <div className="form-group">
                    <label>Статус</label>
                    <select
                      value={newTask.status}
                      onChange={(e) => setNewTask({...newTask, status: e.target.value})}
                    >
                      <option value="pending">⏳ Ожидает</option>
                      <option value="in_progress">⚡ В работе</option>
                      <option value="completed">✅ Завершена</option>
                    </select>
                  </div>
                  <div className="form-group">
                    <label>Приоритет</label>
                    <select
                      value={newTask.priority}
                      onChange={(e) => setNewTask({...newTask, priority: e.target.value})}
                    >
                      <option value="low">🟢 Низкий</option>
                      <option value="medium">🟡 Средний</option>
                      <option value="high">🔴 Высокий</option>
                    </select>
                  </div>
                  {canAssignToUsers && (
                    <div className="form-group">
                      <label>Назначить пользователю</label>
                      <select
                        value={newTask.assigned_to}
                        onChange={(e) => setNewTask({...newTask, assigned_to: e.target.value})}
                      >
                        <option value="">👤 Выберите пользователя</option>
                        {users.map(user => (
                          <option key={user.id} value={user.id}>
                            {user.username} {user.role === 'admin' ? '(🔧 Админ)' : user.role === 'manager' ? '(👔 Менеджер)' : '(👤 Пользователь)'}
                          </option>
                        ))}
                      </select>
                    </div>
                  )}
                </div>
                <div className="form-actions">
                  <button type="submit" className="save-btn">💾 Сохранить</button>
                  <button type="button" className="cancel-btn" onClick={() => setShowTaskForm(false)}>
                    ❌ Отмена
                  </button>
                </div>
              </form>
            )}

            {editingTask && (
              <form className="form" onSubmit={updateTask}>
                <h3>✏️ Редактирование задачи</h3>
                <div className="form-grid">
                  <div className="form-group full-width">
                    <label>Название задачи *</label>
                    <input
                      type="text"
                      value={editingTask.title}
                      onChange={(e) => setEditingTask({...editingTask, title: e.target.value})}
                      required
                    />
                  </div>
                  <div className="form-group full-width">
                    <label>Описание</label>
                    <textarea
                      value={editingTask.description || ''}
                      onChange={(e) => setEditingTask({...editingTask, description: e.target.value})}
                      rows="3"
                    />
                  </div>
                  <div className="form-group">
                    <label>Статус</label>
                    <select
                      value={editingTask.status}
                      onChange={(e) => setEditingTask({...editingTask, status: e.target.value})}
                    >
                      <option value="pending">⏳ Ожидает</option>
                      <option value="in_progress">⚡ В работе</option>
                      <option value="completed">✅ Завершена</option>
                    </select>
                  </div>
                  <div className="form-group">
                    <label>Приоритет</label>
                    <select
                      value={editingTask.priority}
                      onChange={(e) => setEditingTask({...editingTask, priority: e.target.value})}
                    >
                      <option value="low">🟢 Низкий</option>
                      <option value="medium">🟡 Средний</option>
                      <option value="high">🔴 Высокий</option>
                    </select>
                  </div>
                  {canAssignToUsers && (
                    <div className="form-group">
                      <label>Назначить пользователю</label>
                      <select
                        value={editingTask.assigned_to}
                        onChange={(e) => setEditingTask({...editingTask, assigned_to: e.target.value})}
                      >
                        <option value="">👤 Выберите пользователя</option>
                        {users.map(user => (
                          <option key={user.id} value={user.id}>
                            {user.username} {user.role === 'admin' ? '(🔧 Админ)' : user.role === 'manager' ? '(👔 Менеджер)' : '(👤 Пользователь)'}
                          </option>
                        ))}
                      </select>
                    </div>
                  )}
                </div>
                <div className="form-actions">
                  <button type="submit" className="save-btn">💾 Сохранить</button>
                  <button type="button" className="cancel-btn" onClick={() => setEditingTask(null)}>
                    ❌ Отмена
                  </button>
                </div>
              </form>
            )}

            <div className="tasks-list">
              {tasks.length === 0 ? (
                <div className="empty-state">
                  ✅ Нет задач
                </div>
              ) : (
                <div className="cards-grid">
                  {tasks.map(task => (
                    <div key={task.id} className="task-card">
                      <div className="card-header">
                        <div className="task-title">{task.title}</div>
                        <div className="card-actions">
                          <button 
                            className="edit-btn"
                            onClick={() => openEditTask(task)}
                            title="Редактировать"
                          >
                            ✏️
                          </button>
                          <button 
                            className="delete-btn"
                            onClick={() => deleteTask(task.id)}
                            title="Удалить"
                          >
                            🗑️
                          </button>
                        </div>
                      </div>
                      <div className="card-content">
                        {task.description && (
                          <div className="task-description">{task.description}</div>
                        )}
                        <div className="task-assignee">
                          <small>👤 Назначена: {getUserName(task.assigned_to)}</small>
                        </div>
                        <div className="task-meta">
                          <span className={`status-badge status-${task.status}`}>
                            {task.status === 'pending' ? '⏳ Ожидает' : 
                             task.status === 'in_progress' ? '⚡ В работе' : '✅ Завершена'}
                          </span>
                          <span className={`priority-badge priority-${task.priority}`}>
                            {task.priority === 'high' ? '🔴 Высокий' : 
                             task.priority === 'medium' ? '🟡 Средний' : '🟢 Низкий'}
                          </span>
                        </div>
                        <div className="card-footer">
                          <small>Создана: {new Date(task.created_at).toLocaleDateString('ru-RU')}</small>
                          {task.updated_at !== task.created_at && (
                            <small>Обновлена: {new Date(task.updated_at).toLocaleDateString('ru-RU')}</small>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}
      </main>
    </div>
  );
}

export default App;