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
          <li
            key={item.id}
            className="list-group-item list-group-item-action d-flex justify-content-between align-items-start"
            style={{ cursor: "pointer" }}
          >
            <Link to={`/deck/${item.id}`}>
              <div className="ms-2 me-auto">
                <div className="fw-bold">{item.name}</div>
                {item.description}
              </div>
            </Link>
          </li>
        ))}
      </ul>
    </>
  );
}

export default DeckList;
