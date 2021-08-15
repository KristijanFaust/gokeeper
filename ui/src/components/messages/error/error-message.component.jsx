import {faExclamationCircle} from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {useState} from 'react';

import './error-message.styles.scss';

const ErrorMessage = ({children}) => {
  const [clicked, setClicked] = useState(false);

  if (clicked) return null

  return (
    <div className='error-message' onClick={() => setClicked(true)}>
      <FontAwesomeIcon className={'error-icon'} icon={faExclamationCircle}/>
      {children.toUpperCase()}
    </div>
  )
};

export default ErrorMessage;