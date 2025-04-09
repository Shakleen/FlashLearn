import NavBar from "../components/NavBar";
import { useParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
function DeleteConfirmationPage() {
  const { id } = useParams();
  const navigate = useNavigate();

  const handleDelete = async () => {
    try {
      const response = await fetch(`http://localhost:8080/deck/${id}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        throw new Error(`Http error! Status: ${response.status}`);
      }

      toast.success("Deck deleted successfully");
      navigate("/");
    } catch (error) {
      toast.error(
        `Error deleting deck: ${
          error instanceof Error ? error.message : "Unknown error"
        }`
      );
    }
  };

  return (
    <>
      <NavBar />
      <div className="card text-center m-4">
        <h5 className="card-header text-bg-danger">Warning</h5>
        <div className="card-body">
          <h5 className="card-title">
            Are you sure you want to delete this deck?
          </h5>
          <p className="card-text">This action cannot be undone.</p>
          <button onClick={handleDelete} className="btn btn-danger m-2">
            Delete
          </button>
          <button
            onClick={() => {
              navigate(-1);
            }}
            className="btn btn-secondary m-2"
          >
            Cancel
          </button>
        </div>
      </div>
    </>
  );
}

export default DeleteConfirmationPage;
