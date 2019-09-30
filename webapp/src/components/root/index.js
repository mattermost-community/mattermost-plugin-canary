import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {checkCanaryCookie} from 'actions';

import Root from './root';

const mapDispatchToProps = (dispatch) => bindActionCreators({
    checkCanaryCookie,
}, dispatch);

export default connect(null, mapDispatchToProps)(Root);
