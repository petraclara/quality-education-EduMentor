# 🎓 MentorConnect

![MentorConnect](https://mentorconnect-uqsr.onrender.com/)

MentorConnect is an intelligent platform designed to bridge the gap between eager learners and experienced mentors. Users can create specialized profiles, discover mathematically-matched mentors based on skills and goals, and schedule interactive meetings all within a unified platform.

**🌐 Live Demo:** [https://mentorconnect-uqsr.onrender.com](https://mentorconnect-uqsr.onrender.com)
*(Note: As this project is hosted on Render's free tier, the server will spool down after 15 minutes of inactivity. Please allow 30-50 seconds for the backend to wake up on the first load!)*

## ✨ Key Features

- **Smart Matching Algorithm:** Ranks mentors based on overlapping technical skills and specific learning interests using Jaccard Similarity formulas.
- **End-to-End Scheduling Flow:** Learners can request mentorship. If accepted, Mentors can propose multiple date/time slots, meeting types (Zoom, Google Meet, In-Person), and links inline. Learners simply choose the time that works best for them to dynamically finalize a meeting.
- **Role-Based Workflows:** Distinct UI dashboards and backend verification APIs designed specifically for the unique needs of both `Mentors` and `Learners`.
- **Dynamic Database Translation:** Built to operate elegantly using lightweight **SQLite** locally for seamless development and automatically translates schema commands into robust **PostgreSQL** syntax when connected in the cloud.

## 🛠️ Tech Stack

### Frontend
- **Framework:** React + Vite
- **Styling:** Vanilla CSS (Zero generic UI frameworks, utilizing custom Glassmorphism/Modern aesthetic patterns).
- **Routing:** React Router DOM

### Backend
- **Language:** Go (Golang)
- **Database:** SQLite (dev) / PostgreSQL (prod)
- **Authentication:** Custom JWT-based cookie sessions via Go Middleware + Bcrypt password hashing.
- **Configuration:** Cloud-native architecture mapping both the React app and Go API to a single deployed process to minimize server overhead.

---

## ⚡ Local Development

Running MentorConnect locally is incredibly simple. By default, the app initializes its own local `.db` SQLite file so you do not need to install complex local databases or run Docker.

### Prerequisites
- [Go (1.20+)](https://go.dev/)
- [Node.js (18+)](https://nodejs.org/)

### Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/quality-education-EduMentor.git
   cd quality-education-EduMentor
   ```

2. **Start the Go Backend:**
   The backend serves the API on `http://localhost:8080`. When spun up, it will automatically run migrations and optionally seed the database with **8 Demo Mentors** to play with.
   ```bash
   go run main.go
   ```

3. **Start the React Frontend:**
   In a new terminal window, navigate into the frontend module to launch the Vite hot-reloading dev server on `http://localhost:5173`.
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

### Default Seeded Users
The application will automatically inject 8 generated Mentors into the local database upon first boot. You can immediately log into the top-ranked seeded mentor using:
- **Email:** `amara@mentorconnect.io`
- **Password:** `mentor123`

---

## 🚀 Deployment

The project is optimized to run as a **Single Web Service** on platforms like Render to conserve active free-tier limits. 

1. Ensure your frontend dependencies reflect production by replacing `vite.config.js` proxy paths with relative URLs (`/api`).
2. Supply a standard `DATABASE_URL` environment variable connected to a Postgres service (e.g. Neon.tech).
3. The platform will execute `build.sh` automatically and compile both modules into a singular Go executable native to your cloud server!

