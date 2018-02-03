package main

const fastmapQueryf = `
with params as (
    select
        %s :: float as minlon,
        %s :: float as minlat,
        %s :: float as maxlon,
        %s :: float as maxlat
), direct_nodes as (
    select get_nodes_by_box(minlon, minlat, maxlon, maxlat) node
    from params
), all_request_ways as (
    select distinct c.way_id as id
    from direct_nodes n
        join current_way_nodes c on (c.node_id = (n.node).id)
), all_request_nodes as (
    select distinct on (c.node_id) get_node_by_id(c.node_id) as node
    from
        all_request_ways w
        join current_way_nodes c on (c.way_id = w.id)
    where c.node_id not in (
        select (n.node).id
        from direct_nodes n
    )
    union all
    select *
    from direct_nodes n
), relations_from_ways_and_nodes as (
    select distinct m.relation_id as id
    from
        (
            select
                id,
                'Way' :: nwr_enum as type
            from all_request_ways
            union all
            select
                (n.node).id,
                'Node' :: nwr_enum as type
            from all_request_nodes n
        ) wn
        join current_relation_members m on (wn.id = m.member_id and wn.type = m.member_type)
), all_request_relations as (
    select r.id
    from relations_from_ways_and_nodes r
    union
    select rm.relation_id
    from relations_from_ways_and_nodes r2
        join current_relation_members rm on (r2.id = rm.member_id and rm.member_type = 'Relation')
)

select array_to_string( array_agg(osm),'' ) FROM (
    select line
    from (
        -- XML header
        select
            '<?xml version="1.0" encoding="UTF-8"?><osm version="0.6" generator="OpenStreetMap server" copyright="OpenStreetMap and contributors" attribution="http://www.openstreetmap.org/copyright" license="http://opendatacommons.org/licenses/odbl/1-0/">' :: text as line
        union all
        -- bounds header
        select xmlelement(name bounds, xmlattributes (minlat, minlon, maxlat, maxlon)) :: text as line
    from params p
    union all
    -- nodes
    select line :: text
    from (
             select n.node :: xml as line
             from all_request_nodes n
             order by (n.node).id
         ) nodes
    union all
    -- ways
    select line :: text
    from (
             select get_way_by_id(w.id) :: xml as line
             from all_request_ways w
             order by w.id
         ) ways
    union all
    -- relations
    select line :: text
    from
        (
            select get_relation_by_id(r.id) :: xml as line
            from all_request_relations r
            order by r.id
        ) relations
    union all
    -- XML footer
    select '</osm>'
    ) response
) as osm
`
