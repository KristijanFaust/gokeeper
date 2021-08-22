import Dashboard from '../../components/dashboard/dashboard.component';

import './dashboard-page.styles.scss'

const DashboardPage = ({authenticationExpiredCallback}) => {
  return (
    <div className='dashboard-page'>
      <Dashboard authenticationExpiredCallback={authenticationExpiredCallback} />
    </div>
  );
};

export default DashboardPage;