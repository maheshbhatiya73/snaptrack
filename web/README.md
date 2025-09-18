# SnapStack Backup System

A modern, Vercel-style login system with role-based authentication and backup management features.

## Features

- 🎨 **Modern UI**: Vercel-inspired design with glassmorphism effects
- 🔐 **Role-based Authentication**: Standard user, Admin, and Sudo access levels
- 📱 **Responsive Design**: Works perfectly on desktop and mobile devices
- ⚡ **Real-time Status**: Live backup system monitoring
- 🛡️ **Security**: Token-based authentication with expiration
- 🎯 **User-friendly**: Intuitive interface with helpful demo credentials

## Demo Credentials

### Standard User
- **Email**: `user@snapstack.com`
- **Password**: `UserPass2025!`
- **Access**: Basic dashboard and backup viewing

### Admin User
- **Email**: `admin@snapstack.com`
- **Password**: `SnapStack2025!`
- **Access**: Admin panel, user management, system settings

### Sudo User (Root Access)
- **Email**: `admin@snapstack.com`
- **Password**: `sudo123!`
- **Access**: Full system control, emergency operations, root-level backup management

## Getting Started

1. **Install Dependencies**
   ```bash
   npm install
   ```

2. **Start Development Server**
   ```bash
   npm run dev
   ```

3. **Open in Browser**
   Navigate to `http://localhost:3000`

## Project Structure

```
├── components/
│   ├── SnapStackLogo.vue      # Brand logo component
│   ├── DemoCredentials.vue    # Demo credentials helper
│   └── Button.vue            # Reusable button component
├── pages/
│   ├── auth/
│   │   └── login.vue         # Main login page
│   ├── admin/
│   │   ├── index.vue         # Admin dashboard
│   │   └── sudo.vue          # Sudo control panel
│   ├── dashboard.vue         # User dashboard
│   └── index.vue             # Landing page
├── middleware/
│   └── auth.global.js        # Authentication middleware
└── layouts/
    └── default.vue           # Default layout
```

## Key Features Explained

### Authentication System
- **Token-based**: Secure authentication with 24-hour token expiration
- **Role-based Access Control**: Different access levels for different user types
- **Middleware Protection**: Automatic route protection based on user roles

### UI/UX Features
- **Glassmorphism Design**: Modern frosted glass effects
- **Smooth Animations**: Hover effects and transitions
- **Loading States**: Visual feedback during authentication
- **Error Handling**: User-friendly error messages
- **Password Visibility Toggle**: Show/hide password functionality

### Backup System Integration
- **Real-time Status**: Live system monitoring
- **Backup Statistics**: Storage usage, backup counts, and performance metrics
- **Admin Controls**: Full system management capabilities
- **Sudo Operations**: Emergency and root-level operations

## Security Considerations

⚠️ **Important**: This is a demo application. In production:

1. **Replace Static Credentials**: Use a secure backend authentication system
2. **Use HTTP-Only Cookies**: Store tokens in secure, HTTP-only cookies
3. **Implement Proper Hashing**: Use bcrypt or similar for password hashing
4. **Add Rate Limiting**: Prevent brute force attacks
5. **Use HTTPS**: Always use secure connections in production
6. **Validate Input**: Implement proper input validation and sanitization

## Customization

### Branding
- Update the `SnapStackLogo.vue` component to use your own logo
- Modify colors in the Tailwind CSS classes
- Update the company name and branding throughout the application

### Authentication
- Replace the static credentials with your backend API
- Modify the `middleware/auth.global.js` for your authentication flow
- Update the login logic in `pages/auth/login.vue`

### Styling
- Customize the color scheme in the Tailwind classes
- Modify the glassmorphism effects in the login page
- Update the component styling to match your brand

## Technologies Used

- **Nuxt 4**: Vue.js framework with SSR capabilities
- **Tailwind CSS**: Utility-first CSS framework
- **Vue 3**: Progressive JavaScript framework
- **Vue Router**: Client-side routing

## License

This project is for demonstration purposes. Please ensure you implement proper security measures before using in production.

---

**SnapStack Backup System** - Secure, modern, and user-friendly backup management.