import {faSpinner} from '@fortawesome/free-solid-svg-icons';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';

import './spinner.component.scss'

const Spinner = () => (
  <FontAwesomeIcon className='spinner' icon={faSpinner} />
)

export default Spinner;