import { Link } from "react-router-dom";
import { ReactNode } from "react";

interface NavBarProps {
  rightComponents?: ReactNode[];
}

function NavBar({ rightComponents = [] }: NavBarProps) {
  return (
    <nav className="navbar bg-primary" data-bs-theme="dark">
      <div className="container-fluid">
        <Link to={"/"} className="navbar-brand text-light">
          Flash Learn
        </Link>
        <div className="d-flex align-items-center gap-2">
          {rightComponents.map((component, index) => (
            <div key={index} className="text-light">
              {component}
            </div>
          ))}
        </div>
      </div>
    </nav>
  );
}

export default NavBar;
