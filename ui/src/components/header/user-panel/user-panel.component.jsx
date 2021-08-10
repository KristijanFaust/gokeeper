import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faUserCircle} from "@fortawesome/free-solid-svg-icons";

import './user-panel.styles.scss'

const UserPanel = () => {
  let username = localStorage.getItem('username');
  let userPanel = username ? <FontAwesomeIcon icon={faUserCircle}/> : '';

  return (
    <div className='user-panel'>
      {userPanel}
    </div>
  );
};

export default UserPanel;