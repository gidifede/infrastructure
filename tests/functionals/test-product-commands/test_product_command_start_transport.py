from commands.product import *
from queries.product import ProductStatus, ProductDetails, ProductLocation, ProductTrack
from tests.functionals.plan import Plan
from commands.models import Receiver, Sender, Location, Product, Empty, Carrier, FromTo
import time
import utils
import pytest

def arrange_start_transport_cmd_fields(product_id):
    return Product(product_id, 'Poste Delivery Business', 'BOX', []), \
        Carrier('motorbike', 'Postman12345', 'FE764OU'), \
        FromTo('CS', 'Via degli Ulivi 30', '07623', 'Milano', 'Italy'), \
        FromTo('CD', 'Via Ripida 27', '07628', 'Milano Centro', 'Italy')


@pytest.mark.functional
@pytest.mark.skip(reason="story LOI-504")
def test_product_command_start_transport_OK(config,product_id):
    ## Arrange
    product, carrier, from_, to = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier=carrier,
            attributes=[],
            from_=from_,
            to=to
        )
    ]
    
    queries = [
        ProductStatus(
            id = product_id
        ), 
        ProductDetails(
            id = product_id
        ), 
        ProductLocation(
            id = product_id
        ), 
        ProductTrack(
            id = product_id)
    ]
    
    expStatus = 'NetworkTransit'

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=queries,
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()
    assert responses[0].status_code == 202
    time.sleep(5)
    results =  plan.start_query()
    
    ## Assert
    assert len(results) == 4
    # status check
    assert results[0]['product']['id'] == product_id
    assert results[0]['status'] == expStatus
    
    # details check
    assert results[1]['product'] == cmds[0].payload['product']
    assert results[1]['sender'] == {}
    assert results[1]['receiver'] == {}

    # location check
    assert results[2]['product'] == cmds[0].payload['product']
    assert results[2]['location'] == {}

    # track check
    assert results[3]['product']['id'] == product_id
    assert len(results[3]['trackList']) == 1
    assert results[3]['trackList'][0]['description'] == expStatus

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_missing_product_KO(config,product_id):
    ## Arrange
    _, carrier, from_, to = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                from_=from_,
                to=to
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['product']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_missing_carrier_KO(config,product_id):
    ## Arrange
    product, _, from_, to = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [ 
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                attributes=[],
                from_=from_,
                to=to
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['carrier']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_missing_timestamp_KO(config,product_id):
    ## Arrange
    product, carrier, from_, to = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
                product=product,
                carrier=carrier,
                attributes=[],
                from_=from_,
                to=to
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['timestamp']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_bad_timestamp_KO(config,product_id):
    ## Arrange
    product, carrier, from_, to = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
                product=product,
                timestamp=123456789,
                carrier=carrier,
                attributes=[],
                from_=from_,
                to=to
            ),
        StartProductTransport(
                product=product,
                timestamp='A bad string',
                carrier=carrier,
                attributes=[]
            ),
        StartProductTransport(
                product=product,
                timestamp='2022-12-07T17:48:33',
                carrier=carrier,
                attributes=[]
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    assert len(responses) == len(cmds)
    for response in responses:
        assert response.status_code == 400

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_missing_attributes_KO(config,product_id):
    ## Arrange
    product, carrier, from_, to = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
                product=product,
                carrier=carrier,
                timestamp=utils.get_timestamp(),
                from_=from_,
                to=to
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['attributes']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_bad_attributes_KO(config,product_id):
    ## Arrange
    product, carrier, from_, to = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
                product=product,
                carrier=carrier,
                timestamp=utils.get_timestamp(),
                attributes='a string',
                from_=from_,
                to=to
            ),
        StartProductTransport(
                product=product,
                carrier=carrier,
                timestamp=utils.get_timestamp(),
                attributes=5,
                from_=from_,
                to=to
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    assert len(responses) == len(cmds)
    for response in responses:
        assert response.status_code == 400

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_missing_from_KO(config,product_id):
    ## Arrange
    product, carrier, _, to = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
                product=product,
                carrier=carrier,
                attributes=[],
                timestamp=utils.get_timestamp(),
                to=to
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['from']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_bad_from_KO(config,product_id):
    ## Arrange
    product, carrier, _, to = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                from_=FromTo([], '', '', '', ''),
                to=to
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                from_=FromTo('', [], '', '', ''),
                to=to
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                from_=FromTo('', '', [], '', ''),
                to=to
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                from_=FromTo('', '', '', [], ''),
                to=to
            ),
        StartProductTransport(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier=carrier,
            attributes=[],
            from_=FromTo('', '', '', '', []),
            to=to
        ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                from_=FromTo(123, '', '', '', ''),
                to=to
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                from_=FromTo('', 123, '', '', ''),
                to=to
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                from_=FromTo('', '', 123, '', ''),
                to=to
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                from_=FromTo('', '', '', 123, ''),
                to=to
            ),
        StartProductTransport(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier=carrier,
            attributes=[],
            from_=FromTo('', '', '', '', 123),
            to=to
        ),
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    assert len(responses) == len(cmds)
    for response in responses:
        assert response.status_code == 400

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_missing_to_KO(config,product_id):
    ## Arrange
    product, carrier, from_, _ = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
                product=product,
                carrier=carrier,
                attributes=[],
                timestamp=utils.get_timestamp(),
                from_=from_
            )
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    try:
        cmds[0].payload['to']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_bad_to_KO(config,product_id):
    ## Arrange
    product, carrier, from_, _ = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                to=FromTo([], '', '', '', ''),
                from_=from_
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                to=FromTo('', [], '', '', ''),
                from_=from_
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                to=FromTo('', '', [], '', ''),
                from_=from_
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                to=FromTo('', '', '', [], ''),
                from_=from_
            ),
        StartProductTransport(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier=carrier,
            attributes=[],
            to=FromTo('', '', '', '', []),
            from_=from_
        ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                to=FromTo(123, '', '', '', ''),
                from_=from_
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                to=FromTo('', 123, '', '', ''),
                from_=from_
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                to=FromTo('', '', 123, '', ''),
                from_=from_
            ),
        StartProductTransport(
                product=product,
                timestamp=utils.get_timestamp(),
                carrier=carrier,
                attributes=[],
                to=FromTo('', '', '', 123, ''),
                from_=from_
            ),
        StartProductTransport(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier=carrier,
            attributes=[],
            to=FromTo('', '', '', '', 123),
            from_=from_
        ),
    ]

    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key=config['API_KEY']
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    assert len(responses) == len(cmds)
    for response in responses:
        assert response.status_code == 400

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_missing_api_key_OK(config,product_id):
    ## Arrange
    product, carrier, from_, to = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier=carrier,
            attributes=[],
            from_=from_,
            to=to
        )
    ]
    
    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        # missing api key
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    for response in responses:
        assert response.status_code == 403

@pytest.mark.functional
#@pytest.mark.skip(reason="story LOI-489")
def test_product_command_start_transport_wrong_api_key_OK(config,product_id):
    ## Arrange
    product, carrier, from_, to = arrange_start_transport_cmd_fields(product_id=product_id)
    cmds = [
        StartProductTransport(
            product=product,
            timestamp=utils.get_timestamp(),
            carrier=carrier,
            attributes=[],
            from_=from_,
            to=to
        )
    ]
    
    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key='fakeAPiKey'
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    for response in responses:
        assert response.status_code == 403
    