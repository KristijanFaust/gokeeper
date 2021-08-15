import {faCheckCircle} from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {useState} from 'react';

import './notification-message.styles.scss';

const NotificationMessage = ({children}) => {
  const [clicked, setClicked] = useState(false);

  if (clicked) return null;

  return (
    <div className='notification-message' onClick={() => setClicked(true)}>
      <FontAwesomeIcon className={'notification-icon'} icon={faCheckCircle}/>
      {children.toUpperCase()}
    </div>
  )
};

export default NotificationMessage;