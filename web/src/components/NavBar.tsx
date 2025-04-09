import { Link } from "react-router-dom";

function NavBar() {
  return (
    <nav className="navbar bg-body-tertiary" data-bs-theme="dark">
      <div className="container-fluid">
        <Link to={"/"} className="navbar-brand">
          Flash Learn
        </Link>
      </div>
    </nav>
  );
}

export default NavBar;
