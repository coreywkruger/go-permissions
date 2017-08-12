INSERT INTO apps (id, name) VALUES 
  ('697d78cb-b56d-41ad-a7a3-e2e08ebb09fb', 'TacoApp');

INSERT INTO roles (id, name, app_id) VALUES 
  ('c51003fc-2ae4-4296-9d5e-325c76a40316', 'admin', '697d78cb-b56d-41ad-a7a3-e2e08ebb09fb'),
  ('c1688c91-b818-4917-a20e-b95a2006c07f', 'customer', '697d78cb-b56d-41ad-a7a3-e2e08ebb09fb');

INSERT INTO permissions (id, name, app_id) VALUES 
  ('5bee1c60-43e4-460e-80ae-b7c3b8774033', 'read', '697d78cb-b56d-41ad-a7a3-e2e08ebb09fb'),
  ('73017965-b16c-4c6e-9ec1-1e1272594648', 'write', '697d78cb-b56d-41ad-a7a3-e2e08ebb09fb'),
  ('28a212cc-51eb-4e17-95e1-2baa65e55b16', 'delete', '697d78cb-b56d-41ad-a7a3-e2e08ebb09fb');

INSERT INTO role_permissions (id, permission_id, role_id) VALUES 
-- admin
  ('87c5d2bd-13d7-447f-ba63-84eeaa0ac928', '5bee1c60-43e4-460e-80ae-b7c3b8774033', 'c51003fc-2ae4-4296-9d5e-325c76a40316'),
  ('d87fea35-4344-4930-9a59-be976df0266a', '73017965-b16c-4c6e-9ec1-1e1272594648', 'c51003fc-2ae4-4296-9d5e-325c76a40316'),
  ('3d144533-fafd-4050-9dc5-7ab3e479cd73', '28a212cc-51eb-4e17-95e1-2baa65e55b16', 'c51003fc-2ae4-4296-9d5e-325c76a40316'),
-- customer
  ('5b4aec52-44c7-4efa-a5b8-dd61b39a1b4f', '5bee1c60-43e4-460e-80ae-b7c3b8774033', 'c1688c91-b818-4917-a20e-b95a2006c07f');

INSERT INTO entity_roles (id, entity_id, role_id) VALUES 
  ('63948426-c016-4fa9-b2ed-e90589a4deb7', '809e5e2f-0555-4d81-8f91-d6d8f0d4ea79', 'c51003fc-2ae4-4296-9d5e-325c76a40316'),
  ('2ff8542c-d34d-491c-a133-238d0bdd12fa', '07df4a77-6243-41cd-a421-90c524ef2203', 'c1688c91-b818-4917-a20e-b95a2006c07f');