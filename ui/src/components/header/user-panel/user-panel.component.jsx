import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faUserCircle} from '@fortawesome/free-solid-svg-icons';
import {useState, useRef} from 'react';

import UserPanelDropdown from './user-panel-dropdown/user-panel-dropdown.component';
import useOutsideClick from '../../../utility/hooks/use-outside-click.hook';

import './user-panel.styles.scss';

const UserPanel = ({signOutCallback}) => {
  const username = localStorage.getItem('username');
  const userPanel = username ? <FontAwesomeIcon className={'user-panel-icon'} icon={faUserCircle}/> : null;

  const [clicked, setClicked] = useState(false);
  const ref = useRef(null);
  useOutsideClick(ref, setClicked, false);

  const userPanelDropdown = clicked ? <UserPanelDropdown signOutCallback={signOutCallback} /> : null;

  return (
    <div className='user-panel' ref={ref} onClick={() => setClicked(true)}>
      <div className='username'>
        {username}
      </div>
      {userPanel}
      {userPanelDropdown}
    </div>
  );
};

export default UserPanel;