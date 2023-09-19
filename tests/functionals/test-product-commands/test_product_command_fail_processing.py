from commands.product import *
from queries.product import ProductStatus, ProductDetails, ProductLocation, ProductTrack
from tests.functionals.plan import Plan
from commands.models import Receiver, Sender, Location, Product, Empty
import time
import utils
import pytest

def arrange_fail_processing_cmd_fields(product_id):
    return Product(product_id, 'Poste Delivery Business', 'BOX', []), \
        Location('UP', 'Roma', 'Viale Europa 190', '00144', 'Italia', '55Y90', [])


@pytest.mark.functional
@pytest.mark.skip(reason="bug LOI-504")
def test_product_command_fail_processing_OK(config,product_id):
    ## Arrange
    product, location = arrange_fail_processing_cmd_fields(product_id=product_id)
    cmds = [ 
        FailProductProcessing(
            product=product,
            timestamp=utils.get_timestamp(),
            location=location,
            notes='Receiver not found',
            attributes=[]
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
    
    expStatus = 'Unworkable'

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
    assert results[2]['location'] == cmds[0].payload['location']

    # track check
    assert results[3]['product']['id'] == product_id
    assert len(results[3]['trackList']) == 1
    assert results[3]['trackList'][0]['description'] == expStatus

@pytest.mark.functional
def test_product_command_fail_processing_missing_product_KO(config,product_id):
    ## Arrange
    _, location = arrange_fail_processing_cmd_fields(product_id=product_id)
    cmds = [
        FailProductProcessing(
                timestamp=utils.get_timestamp(),
                location=location,
                notes='Receiver not found',
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
    try:
        cmds[0].payload['product']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_fail_processing_missing_location_KO(config,product_id):
    ## Arrange
    product, _ = arrange_fail_processing_cmd_fields(product_id=product_id)
    cmds = [
        FailProductProcessing(
                product=product,
                timestamp=utils.get_timestamp(),
                notes='Receiver not found',
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
    try:
        cmds[0].payload['location']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_fail_processing_missing_timestamp_KO(config,product_id):
    ## Arrange
    product, location = arrange_fail_processing_cmd_fields(product_id=product_id)
    cmds = [
        FailProductProcessing(
                product=product,
                location=location,
                notes='Receiver not found',
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
    try:
        cmds[0].payload['timestamp']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_fail_processing_bad_timestamp_KO(config,product_id):
    ## Arrange
    product, location = arrange_fail_processing_cmd_fields(product_id=product_id)
    cmds = [
        FailProductProcessing(
                product=product,
                timestamp=123456789,
                location=location,
                notes='Receiver not found',
                attributes=[]
            ),
        FailProductProcessing(
                product=product,
                timestamp='A bad string',
                location=location,
                notes='Receiver not found',
                attributes=[]
            ),
        FailProductProcessing(
                product=product,
                timestamp='2022-12-07T17:48:33',
                location=location,
                notes='Receiver not found',
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
def test_product_command_fail_processing_missing_notes_KO(config,product_id):
    ## Arrange
    product, location = arrange_fail_processing_cmd_fields(product_id=product_id)
    cmds = [
        FailProductProcessing(
                product=product,
                timestamp=utils.get_timestamp(),
                location=location,
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
    try:
        cmds[0].payload['notes']
        assert False
    except KeyError:
        assert True
    assert len(responses) == 1
    assert responses[0].status_code == 400

@pytest.mark.functional
def test_product_command_fail_processing_bad_notes_KO(config,product_id):
    ## Arrange
    product, location = arrange_fail_processing_cmd_fields(product_id=product_id)
    cmds = [
        FailProductProcessing(
                product=product,
                timestamp=utils.get_timestamp(),
                location=location,
                notes=12345,
                attributes=[]
            ),
        FailProductProcessing(
                product=product,
                timestamp=utils.get_timestamp(),
                location=location,
                notes=[],
                attributes=[]
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
def test_product_command_fail_processing_missing_attributes_KO(config,product_id):
    ## Arrange
    product, location = arrange_fail_processing_cmd_fields(product_id=product_id)
    cmds = [
        FailProductProcessing(
                product=product,
                timestamp=utils.get_timestamp(),
                location=location,
                notes='Missing Receiver'
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
def test_product_command_fail_processing_bad_attributes_KO(config,product_id):
    ## Arrange
    product, location = arrange_fail_processing_cmd_fields(product_id=product_id)
    cmds = [
        FailProductProcessing(
                product=product,
                timestamp=utils.get_timestamp(),
                location=location,
                notes='Missing Receiver',
                attributes='wrong'
            ),
        FailProductProcessing(
                product=product,
                timestamp=utils.get_timestamp(),
                location=location,
                notes='Missing Receiver',
                attributes=12345
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
def test_product_command_fail_processing_missing_api_key_OK(config,product_id):
    ## Arrange
    product, location = arrange_fail_processing_cmd_fields(product_id=product_id)
    cmds = [ 
        FailProductProcessing(
            product=product,
            timestamp=utils.get_timestamp(),
            location=location,
            notes='Receiver not found',
            attributes=[]
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
def test_product_command_fail_processing_wrong_api_key_OK(config,product_id):
    ## Arrange
    product, location = arrange_fail_processing_cmd_fields(product_id=product_id)
    cmds = [ 
        FailProductProcessing(
            product=product,
            timestamp=utils.get_timestamp(),
            location=location,
            notes='Receiver not found',
            attributes=[]
        )
    ]
    
    plan = Plan(
        id=product_id,
        commands=cmds,
        queries=[],
        base_api=config['BASE_API'],
        api_key='fakeAPIKey'
    )
    
    ## Execute
    responses = plan.start_command()

    ## Assert
    for response in responses:
        assert response.status_code == 403