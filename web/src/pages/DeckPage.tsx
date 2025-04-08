import { useParams } from "react-router-dom";

function DeckPage() {
  const { id } = useParams();

  return (
    <>
      <nav className="navbar bg-body-tertiary" data-bs-theme="dark">
        <div className="container-fluid">
          <a className="navbar-brand" href="#">
            Navbar
          </a>
        </div>
      </nav>
      <h1>This is the deck {id} page</h1>
    </>
  );
}

export default DeckPage;
