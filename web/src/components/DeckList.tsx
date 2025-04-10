import { Link } from "react-router-dom";

export interface DeckItem {
  id: string | number;
  name: string;
  description: string;
}

interface DeckListProps {
  items: DeckItem[];
}

function DeckList(props: DeckListProps) {
  return (
    <>
      <ul className="list-group">
        {props.items.map((item) => (
          <Link
            to={`/deck/${item.id}`}
            state={{ deck: item }}
            className="text-decoration-none"
          >
            <li
              key={item.id}
              className="list-group-item list-group-item-action d-flex justify-content-between align-items-start"
              style={{ cursor: "pointer" }}
            >
              <div className="ms-2 me-auto text-dark">
                <div className="fw-bold fs-4">{item.name}</div>
                {item.description}
              </div>
            </li>
          </Link>
        ))}
      </ul>
    </>
  );
}

export default DeckList;
