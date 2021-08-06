import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faUserCircle} from "@fortawesome/free-solid-svg-icons";

import './user-panel.styles.scss'

const UserPanel = () => (
  <div className='user-panel'>
    <FontAwesomeIcon icon={faUserCircle} />
  </div>
);

export default UserPanel;