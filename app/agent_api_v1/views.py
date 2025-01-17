from flask import jsonify, request, current_app,abort
from . import api
from app.models import *
from flask_login import login_required,current_user
from app.utils.decorators import roles_required, has_valid_tenant_key, has_valid_tenant_and_agent
from app.utils.misc import handle_collection
import arrow
import json
import zlib

from app.logging import logger

@api.route('/register', methods=['POST'])
@has_valid_tenant_key
def register_agent(tenant):
    logger.info("[/register] Entered register agent endpoint")

    if agent_key := request.headers.get("aid"):
        logger.info(f"[/register] Agent key found in headers: {agent_key}")
        if agent := Agent.find(agent_key):
            logger.info(f"[/register] Agent already registered with key: {agent_key}")
            return jsonify({"registered":True,"msg":"agent is already registered","agent_id":agent.key})

    payload = request.get_json()
    logger.info("[/register] PAYLOAD: {}".format(payload))

    agent_id = tenant.register_agent(payload)
    logger.info("[/register] AGENT_ID: {}".format(agent_id))

    return jsonify({"registered":True,"msg":"agent successfully registered","agent_id":agent_id})

@api.route("/policy",methods=["GET","POST"])
@has_valid_tenant_and_agent
def get_policy(agent):
    logger.info("[/policy] Entered get_policy endpoint")
    data = request.get_json()
    if data:
        for key, value in data.items():
            setattr(agent, key, value)
    agent.last_active = arrow.utcnow().datetime
    db.session.commit()
    return jsonify(agent.policy_format())

@api.route("/compliance",methods=["POST"])
@has_valid_tenant_and_agent
def post_compliance(agent):
    logger.info("[/compliance] Entered post_compliance endpoint")
    data = request.get_json()
    task = ComplianceTask.query.filter(ComplianceTask.name == data["name"]).first()
    record = {
        "validate": data["validate"]["output"],
        "validate_rtn_code": int(data["validate"]["rtn_code"]),
        "agent_id": agent.id,
        "tenant_id": agent.tenant_id
    }
    if data.get("enforce"):
        record["enforce"] = data["enforce"]["output"]
        record["enforce_rtn_code"] = int(data["enforce"]["rtn_code"])
    result = ComplianceTaskResults(**record)
    task.results.append(result)
    db.session.commit()
    return jsonify({"message":"ok"})

@api.route("/collection",methods=["POST"])
@has_valid_tenant_and_agent
def data_collection(agent):
    logger.info("[/collection] Entered data_collection endpoint")
    data = zlib.decompress(request.data).decode("utf-8")
    record = json.loads(data)
    handle_collection(agent, record)
    agent.last_active = arrow.utcnow().datetime
    return jsonify({"message":"ok"})
